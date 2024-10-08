// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source client.go -destination client_generated.go -package client . Client
//

// Package client is a generated GoMock package.
package client

import (
	context "context"
	json "encoding/json"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreatePayout mocks base method.
func (m *MockClient) CreatePayout(ctx context.Context, quote Quote, targetAccount uint64, transactionID string) (*Payout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayout", ctx, quote, targetAccount, transactionID)
	ret0, _ := ret[0].(*Payout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePayout indicates an expected call of CreatePayout.
func (mr *MockClientMockRecorder) CreatePayout(ctx, quote, targetAccount, transactionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayout", reflect.TypeOf((*MockClient)(nil).CreatePayout), ctx, quote, targetAccount, transactionID)
}

// CreateQuote mocks base method.
func (m *MockClient) CreateQuote(ctx context.Context, profileID, currency string, amount json.Number) (Quote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuote", ctx, profileID, currency, amount)
	ret0, _ := ret[0].(Quote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateQuote indicates an expected call of CreateQuote.
func (mr *MockClientMockRecorder) CreateQuote(ctx, profileID, currency, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuote", reflect.TypeOf((*MockClient)(nil).CreateQuote), ctx, profileID, currency, amount)
}

// CreateTransfer mocks base method.
func (m *MockClient) CreateTransfer(ctx context.Context, quote Quote, targetAccount uint64, transactionID string) (*Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransfer", ctx, quote, targetAccount, transactionID)
	ret0, _ := ret[0].(*Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransfer indicates an expected call of CreateTransfer.
func (mr *MockClientMockRecorder) CreateTransfer(ctx, quote, targetAccount, transactionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransfer", reflect.TypeOf((*MockClient)(nil).CreateTransfer), ctx, quote, targetAccount, transactionID)
}

// CreateWebhook mocks base method.
func (m *MockClient) CreateWebhook(ctx context.Context, profileID uint64, name, triggerOn, url, version string) (*webhookSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebhook", ctx, profileID, name, triggerOn, url, version)
	ret0, _ := ret[0].(*webhookSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebhook indicates an expected call of CreateWebhook.
func (mr *MockClientMockRecorder) CreateWebhook(ctx, profileID, name, triggerOn, url, version any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebhook", reflect.TypeOf((*MockClient)(nil).CreateWebhook), ctx, profileID, name, triggerOn, url, version)
}

// DeleteWebhooks mocks base method.
func (m *MockClient) DeleteWebhooks(ctx context.Context, profileID uint64, subscriptionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebhooks", ctx, profileID, subscriptionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebhooks indicates an expected call of DeleteWebhooks.
func (mr *MockClientMockRecorder) DeleteWebhooks(ctx, profileID, subscriptionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebhooks", reflect.TypeOf((*MockClient)(nil).DeleteWebhooks), ctx, profileID, subscriptionID)
}

// GetBalance mocks base method.
func (m *MockClient) GetBalance(ctx context.Context, profileID, balanceID uint64) (*Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, profileID, balanceID)
	ret0, _ := ret[0].(*Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockClientMockRecorder) GetBalance(ctx, profileID, balanceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockClient)(nil).GetBalance), ctx, profileID, balanceID)
}

// GetBalances mocks base method.
func (m *MockClient) GetBalances(ctx context.Context, profileID uint64) ([]Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalances", ctx, profileID)
	ret0, _ := ret[0].([]Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalances indicates an expected call of GetBalances.
func (mr *MockClientMockRecorder) GetBalances(ctx, profileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalances", reflect.TypeOf((*MockClient)(nil).GetBalances), ctx, profileID)
}

// GetPayout mocks base method.
func (m *MockClient) GetPayout(ctx context.Context, payoutID string) (*Payout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayout", ctx, payoutID)
	ret0, _ := ret[0].(*Payout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPayout indicates an expected call of GetPayout.
func (mr *MockClientMockRecorder) GetPayout(ctx, payoutID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayout", reflect.TypeOf((*MockClient)(nil).GetPayout), ctx, payoutID)
}

// GetProfiles mocks base method.
func (m *MockClient) GetProfiles(ctx context.Context) ([]Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfiles", ctx)
	ret0, _ := ret[0].([]Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfiles indicates an expected call of GetProfiles.
func (mr *MockClientMockRecorder) GetProfiles(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfiles", reflect.TypeOf((*MockClient)(nil).GetProfiles), ctx)
}

// GetRecipientAccount mocks base method.
func (m *MockClient) GetRecipientAccount(ctx context.Context, accountID uint64) (*RecipientAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecipientAccount", ctx, accountID)
	ret0, _ := ret[0].(*RecipientAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecipientAccount indicates an expected call of GetRecipientAccount.
func (mr *MockClientMockRecorder) GetRecipientAccount(ctx, accountID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecipientAccount", reflect.TypeOf((*MockClient)(nil).GetRecipientAccount), ctx, accountID)
}

// GetRecipientAccounts mocks base method.
func (m *MockClient) GetRecipientAccounts(ctx context.Context, profileID uint64, pageSize int, seekPositionForNext uint64) (*RecipientAccountsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecipientAccounts", ctx, profileID, pageSize, seekPositionForNext)
	ret0, _ := ret[0].(*RecipientAccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecipientAccounts indicates an expected call of GetRecipientAccounts.
func (mr *MockClientMockRecorder) GetRecipientAccounts(ctx, profileID, pageSize, seekPositionForNext any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecipientAccounts", reflect.TypeOf((*MockClient)(nil).GetRecipientAccounts), ctx, profileID, pageSize, seekPositionForNext)
}

// GetTransfer mocks base method.
func (m *MockClient) GetTransfer(ctx context.Context, transferID string) (*Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransfer", ctx, transferID)
	ret0, _ := ret[0].(*Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransfer indicates an expected call of GetTransfer.
func (mr *MockClientMockRecorder) GetTransfer(ctx, transferID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransfer", reflect.TypeOf((*MockClient)(nil).GetTransfer), ctx, transferID)
}

// GetTransfers mocks base method.
func (m *MockClient) GetTransfers(ctx context.Context, profileID uint64, offset, limit int) ([]Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransfers", ctx, profileID, offset, limit)
	ret0, _ := ret[0].([]Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransfers indicates an expected call of GetTransfers.
func (mr *MockClientMockRecorder) GetTransfers(ctx, profileID, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransfers", reflect.TypeOf((*MockClient)(nil).GetTransfers), ctx, profileID, offset, limit)
}

// ListWebhooksSubscription mocks base method.
func (m *MockClient) ListWebhooksSubscription(ctx context.Context, profileID uint64) ([]webhookSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWebhooksSubscription", ctx, profileID)
	ret0, _ := ret[0].([]webhookSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebhooksSubscription indicates an expected call of ListWebhooksSubscription.
func (mr *MockClientMockRecorder) ListWebhooksSubscription(ctx, profileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebhooksSubscription", reflect.TypeOf((*MockClient)(nil).ListWebhooksSubscription), ctx, profileID)
}

// TranslateBalanceUpdateWebhook mocks base method.
func (m *MockClient) TranslateBalanceUpdateWebhook(ctx context.Context, payload []byte) (balanceUpdateWebhookPayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TranslateBalanceUpdateWebhook", ctx, payload)
	ret0, _ := ret[0].(balanceUpdateWebhookPayload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TranslateBalanceUpdateWebhook indicates an expected call of TranslateBalanceUpdateWebhook.
func (mr *MockClientMockRecorder) TranslateBalanceUpdateWebhook(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TranslateBalanceUpdateWebhook", reflect.TypeOf((*MockClient)(nil).TranslateBalanceUpdateWebhook), ctx, payload)
}

// TranslateTransferStateChangedWebhook mocks base method.
func (m *MockClient) TranslateTransferStateChangedWebhook(ctx context.Context, payload []byte) (Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TranslateTransferStateChangedWebhook", ctx, payload)
	ret0, _ := ret[0].(Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TranslateTransferStateChangedWebhook indicates an expected call of TranslateTransferStateChangedWebhook.
func (mr *MockClientMockRecorder) TranslateTransferStateChangedWebhook(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TranslateTransferStateChangedWebhook", reflect.TypeOf((*MockClient)(nil).TranslateTransferStateChangedWebhook), ctx, payload)
}
