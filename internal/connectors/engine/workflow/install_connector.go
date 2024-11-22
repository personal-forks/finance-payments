package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type InstallConnector struct {
	ConnectorID models.ConnectorID
	Config      models.Config
	RawConfig   json.RawMessage
}

func (w Workflow) runInstallConnector(
	ctx workflow.Context,
	installConnector InstallConnector,
) error {
	// Second step: install the connector via the plugin and get the list of
	// capabilities and the workflow of polling data
	installResponse, err := activities.PluginInstallConnector(
		// disable retries as grpc plugin boot command cannot be run more than once by the go-plugin client
		// this also causes API install calls to fail immediately which is more desirable in the case that a plugin is timing out or not compiled correctly
		maximumAttemptsRetryContext(ctx, 1),
		installConnector.ConnectorID,
		installConnector.RawConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to install connector")
	}

	// Third step: store the capabilities of the connector
	if err := activities.StorageCapabilitiesStore(
		infiniteRetryContext(ctx),
		installConnector.ConnectorID,
		installResponse.Capabilities,
	); err != nil {
		return errors.Wrap(err, "failed to store capabilities")
	}

	// Fourth step: store the workflow of the connector
	err = activities.StorageConnectorTasksTreeStore(infiniteRetryContext(ctx), installConnector.ConnectorID, installResponse.Workflow)
	if err != nil {
		return errors.Wrap(err, "failed to store tasks tree")
	}

	if len(installResponse.WebhooksConfigs) > 0 {
		configs := make([]models.WebhookConfig, 0, len(installResponse.WebhooksConfigs))
		for _, webhookConfig := range installResponse.WebhooksConfigs {
			configs = append(configs, models.WebhookConfig{
				Name:        webhookConfig.Name,
				ConnectorID: installConnector.ConnectorID,
				URLPath:     webhookConfig.URLPath,
			})
		}

		err = activities.StorageWebhooksConfigsStore(infiniteRetryContext(ctx), configs)
		if err != nil {
			return errors.Wrap(err, "failed to store webhooks configs")
		}
	}

	// Fifth step: launch the workflow tree, do not wait for the result
	// by using the GetChildWorkflowExecution function that returns a future
	// which will be ready when the child workflow has successfully started.
	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(
			ctx,
			workflow.ChildWorkflowOptions{
				WorkflowID:            fmt.Sprintf("run-tasks-%s-%s", w.stack, installConnector.ConnectorID.String()),
				WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
				TaskQueue:             installConnector.ConnectorID.String(),
				ParentClosePolicy:     enums.PARENT_CLOSE_POLICY_ABANDON,
				SearchAttributes: map[string]interface{}{
					SearchAttributeStack: w.stack,
				},
			},
		),
		Run,
		installConnector.Config,
		installConnector.ConnectorID,
		nil,
		[]models.ConnectorTaskTree(installResponse.Workflow),
	).GetChildWorkflowExecution().Get(ctx, nil); err != nil {
		applicationError := &temporal.ApplicationError{}
		if errors.As(err, &applicationError) {
			if applicationError.Type() != "ChildWorkflowExecutionAlreadyStartedError" {
				return err
			}
		} else {
			return errors.Wrap(err, "running next workflow")
		}
	}

	return nil
}

const RunInstallConnector = "InstallConnector"
