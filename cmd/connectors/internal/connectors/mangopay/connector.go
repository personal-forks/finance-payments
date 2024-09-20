package mangopay

import (
	"context"
	"errors"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
)

const name = models.ConnectorProviderMangopay

var (
	mainTaskDescriptor = TaskDescriptor{
		Name: "Main task to periodically fetch users and transactions",
		Key:  taskNameMain,
	}
)

type Connector struct {
	logger logging.Logger
	cfg    Config

	taskMemoryState taskMemoryState
}

func newConnector(logger logging.Logger, cfg Config) *Connector {
	return &Connector{
		logger: logger.WithFields(map[string]any{
			"component": "connector",
		}),
		cfg: cfg,
	}
}

func (c *Connector) UpdateConfig(ctx task.ConnectorContext, config models.ConnectorConfigObject) error {
	cfg, ok := config.(Config)
	if !ok {
		return connectors.ErrInvalidConfig
	}

	restartTask := c.cfg.PollingPeriod.Duration != cfg.PollingPeriod.Duration

	c.cfg = cfg

	if restartTask {
		taskDescriptor, err := models.EncodeTaskDescriptor(mainTaskDescriptor)
		if err != nil {
			return err
		}

		return ctx.Scheduler().Schedule(ctx.Context(), taskDescriptor, models.TaskSchedulerOptions{
			// We want to polling every c.cfg.PollingPeriod.Duration seconds the users
			// and their transactions.
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       c.cfg.PollingPeriod.Duration,
			// No need to restart this task, since the connector is not existing or
			// was uninstalled previously, the task does not exists in the database
			RestartOption: models.OPTIONS_STOP_AND_RESTART,
		})
	}

	return nil
}

func (c *Connector) Install(ctx task.ConnectorContext) error {
	// Create hooks on mangopay sync
	createWebhookDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: "First task to create webhooks",
		Key:  taskNameCreateWebhook,
	})
	if err != nil {
		return err
	}

	if err := ctx.Scheduler().Schedule(ctx.Context(), createWebhookDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW_SYNC,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	}); err != nil {
		return err
	}

	mainDescriptor, err := models.EncodeTaskDescriptor(mainTaskDescriptor)
	if err != nil {
		return err
	}

	return ctx.Scheduler().Schedule(ctx.Context(), mainDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
		Duration:       c.cfg.PollingPeriod.Duration,
		// No need to restart this task, since the connector is not existing or
		// was uninstalled previously, the task does not exists in the database
		RestartOption: models.OPTIONS_RESTART_NEVER,
	})
}

func (c *Connector) Uninstall(ctx context.Context) error {
	return nil
}

func (c *Connector) Resolve(descriptor models.TaskDescriptor) task.Task {
	taskDescriptor, err := models.DecodeTaskDescriptor[TaskDescriptor](descriptor)
	if err != nil {
		panic(err)
	}

	return resolveTasks(c.logger, c.cfg, &c.taskMemoryState)(taskDescriptor)
}

func (c *Connector) SupportedCurrenciesAndDecimals() map[string]int {
	return supportedCurrenciesWithDecimal
}

func (c *Connector) InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error {
	taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:       "Initiate payment",
		Key:        taskNameInitiatePayment,
		TransferID: transfer.ID.String(),
	})
	if err != nil {
		return err
	}

	scheduleOption := models.OPTIONS_RUN_NOW_SYNC
	scheduledAt := transfer.ScheduledAt
	if !scheduledAt.IsZero() {
		scheduleOption = models.OPTIONS_RUN_SCHEDULED_AT
	}

	err = ctx.Scheduler().Schedule(ctx.Context(), taskDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: scheduleOption,
		ScheduleAt:     scheduledAt,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	})
	if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
		return err
	}

	return nil
}

func (c *Connector) ReversePayment(ctx task.ConnectorContext, transferReversal *models.TransferReversal) error {
	return connectors.ErrNotImplemented
}

func (c *Connector) CreateExternalBankAccount(ctx task.ConnectorContext, bankAccount *models.BankAccount) error {
	descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:          "Create external bank account",
		Key:           taskNameCreateExternalBankAccount,
		BankAccountID: bankAccount.ID,
	})
	if err != nil {
		return err
	}
	if err := ctx.Scheduler().Schedule(ctx.Context(), descriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW_SYNC,
	}); err != nil {
		return err
	}

	return nil
}

var _ connectors.Connector = &Connector{}
