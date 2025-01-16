package workflow

import (
	"context"
	"errors"
	"math/big"

	"github.com/formancehq/go-libs/v2/bun/bunpaginate"
	"github.com/formancehq/go-libs/v2/pointer"
	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (s *UnitTestSuite) Test_ReversePayout_Success() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, req activities.ReversePayoutRequest) (*models.ReversePayoutResponse, error) {
		s.Equal(s.connectorID, req.ConnectorID)
		s.Equal(s.paymentReversalID.Reference, req.Req.PaymentInitiationReversal.Reference)
		return &models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		}, nil
	})
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, payments []models.Payment) error {
		s.Len(payments, 1)
		s.Equal(s.paymentPayoutID, payments[0].ID)
		return nil
	})
	s.env.OnWorkflow(RunSendEvents, mock.Anything, mock.Anything).Once().Return(func(ctx workflow.Context, sendEvents SendEvents) error {
		s.Nil(sendEvents.Balance)
		s.Nil(sendEvents.Account)
		s.Nil(sendEvents.ConnectorReset)
		s.NotNil(sendEvents.Payment)
		s.Nil(sendEvents.PoolsCreation)
		s.Nil(sendEvents.PoolsDeletion)
		s.Nil(sendEvents.BankAccount)
		return nil
	})
	s.env.OnActivity(activities.StoragePaymentInitiationsRelatedPaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, relatedPayment activities.RelatedPayment) error {
		s.Equal(s.paymentInitiationPayout.ID, relatedPayment.PiID)
		s.Equal(s.paymentPayoutID, relatedPayment.PID)
		return nil
	})
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, adj models.PaymentInitiationAdjustment) error {
		s.Equal(models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_REVERSED, adj.Status)
		return nil
	})
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, adj models.PaymentInitiationReversalAdjustment) error {
		s.Equal(models.PAYMENT_INITIATION_REVERSAL_STATUS_PROCESSED, adj.Status)
		return nil
	})
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_SUCCEEDED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.NoError(err)
}

func (s *UnitTestSuite) Test_ReversePayout_PluginReversePayout_Error_Success() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, req activities.ReversePayoutRequest) (*models.ReversePayoutResponse, error) {
		return nil, temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error"))
	})
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, adj models.PaymentInitiationReversalAdjustment) error {
		s.Equal(models.PAYMENT_INITIATION_REVERSAL_STATUS_FAILED, adj.Status)
		return nil
	})
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, adj models.PaymentInitiationAdjustment) error {
		s.Equal(models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_REVERSE_FAILED, adj.Status)
		return nil
	})
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationReversalsGet_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		nil,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationsGet_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		nil,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationAdjustmentsList_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		nil,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_ValidatePaymentInitiationProcessed_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	// No processed payment initiation adjustments
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "PAYMENT_INITIATION_NOT_PROCESSED")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationAdjustmentsList_2_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		nil,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_ValidateReverseAmount_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_REVERSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_REVERSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "CANNOT_REVERSE_MORE_THAN_AMOUNT")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationsAdjusmentsIfStatusEqualStore_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		true,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_ValidateOnlyReverse_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(false, nil)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "ANOTHER_REVERSE_IN_PROGRESS")
}

func (s *UnitTestSuite) Test_ReversePayout_StorageAccountsGet_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		nil,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StorageAccountsGet_2_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		nil,
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentsStore_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(
		&models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_RunSendEvents_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(
		&models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnWorkflow(RunSendEvents, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationsRelatedPaymentsStore_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(
		&models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnWorkflow(RunSendEvents, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsRelatedPaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationsAdjustmentsStore_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(
		&models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnWorkflow(RunSendEvents, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsRelatedPaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationReversalsAdjustmentsStore_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(
		&models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnWorkflow(RunSendEvents, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsRelatedPaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StorageTasksStore_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(
		&models.ReversePayoutResponse{
			Payment: s.pspPaymentReversed,
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnWorkflow(RunSendEvents, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsRelatedPaymentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationReversalsAdjustmentsStore_2_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, req activities.ReversePayoutRequest) (*models.ReversePayoutResponse, error) {
		return nil, temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error"))
	})
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}

func (s *UnitTestSuite) Test_ReversePayout_StoragePaymentInitiationsAdjustmentsStore_2_Error() {
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsGetActivity, mock.Anything, s.paymentReversalID).Once().Return(
		&s.paymentReversal,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsGetActivity, mock.Anything, s.paymentInitiationID).Once().Return(
		&s.paymentInitiationPayout,
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 1,
			HasMore:  false,
			Data: []models.PaymentInitiationAdjustment{
				{
					ID: models.PaymentInitiationAdjustmentID{
						PaymentInitiationID: s.paymentInitiationID,
						CreatedAt:           s.env.Now(),
						Status:              models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					},
					CreatedAt: s.env.Now(),
					Status:    models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_PROCESSED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Error:     nil,
					Metadata: map[string]string{
						"test": "test",
					},
				},
			},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationAdjustmentsListActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(
		&bunpaginate.Cursor[models.PaymentInitiationAdjustment]{
			PageSize: 0,
			HasMore:  false,
			Data:     []models.PaymentInitiationAdjustment{},
		},
		nil,
	)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjusmentsIfStatusEqualStoreActivity, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.SourceAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.StorageAccountsGetActivity, mock.Anything, *s.paymentInitiationPayout.DestinationAccountID).Once().Return(
		&s.account,
		nil,
	)
	s.env.OnActivity(activities.PluginReversePayoutActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, req activities.ReversePayoutRequest) (*models.ReversePayoutResponse, error) {
		return nil, temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error"))
	})
	s.env.OnActivity(activities.StoragePaymentInitiationReversalsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(nil)
	s.env.OnActivity(activities.StoragePaymentInitiationsAdjustmentsStoreActivity, mock.Anything, mock.Anything).Once().Return(
		temporal.NewNonRetryableApplicationError("test", "test", errors.New("test-error")),
	)
	s.env.OnActivity(activities.StorageTasksStoreActivity, mock.Anything, mock.Anything).Once().Return(func(ctx context.Context, task models.Task) error {
		s.Equal(models.TASK_STATUS_FAILED, task.Status)
		return nil
	})

	s.env.ExecuteWorkflow(RunReversePayout, ReversePayout{
		TaskID: models.TaskID{
			Reference:   "test",
			ConnectorID: s.connectorID,
		},
		ConnectorID:                 s.connectorID,
		PaymentInitiationReversalID: s.paymentReversalID,
	})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.ErrorContains(err, "test-error")
}
