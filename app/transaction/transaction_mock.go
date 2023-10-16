// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package transaction is a generated GoMock package.
package transaction

import (
	context "context"
	reflect "reflect"
	model "github.com/jcpribeiro/TransactionApp/model"

	gomock "github.com/golang/mock/gomock"
)

// MockApp is a mock of App interface.
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp.
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance.
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// GetTransactions mocks base method.
func (m *MockApp) GetTransactions(ctx context.Context, transactionIds []string) ([]*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", ctx, transactionIds)
	ret0, _ := ret[0].([]*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions.
func (mr *MockAppMockRecorder) GetTransactions(ctx, transactionIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockApp)(nil).GetTransactions), ctx, transactionIds)
}

// GetTransactionsByPeriod mocks base method.
func (m *MockApp) GetTransactionsByPeriod(ctx context.Context, startDate, endDate string) ([]*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByPeriod", ctx, startDate, endDate)
	ret0, _ := ret[0].([]*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByPeriod indicates an expected call of GetTransactionsByPeriod.
func (mr *MockAppMockRecorder) GetTransactionsByPeriod(ctx, startDate, endDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByPeriod", reflect.TypeOf((*MockApp)(nil).GetTransactionsByPeriod), ctx, startDate, endDate)
}

// GetTransactionsByPeriodEpoch mocks base method.
func (m *MockApp) GetTransactionsByPeriodEpoch(ctx context.Context, startDate, endDate int64) ([]*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByPeriodEpoch", ctx, startDate, endDate)
	ret0, _ := ret[0].([]*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByPeriodEpoch indicates an expected call of GetTransactionsByPeriodEpoch.
func (mr *MockAppMockRecorder) GetTransactionsByPeriodEpoch(ctx, startDate, endDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByPeriodEpoch", reflect.TypeOf((*MockApp)(nil).GetTransactionsByPeriodEpoch), ctx, startDate, endDate)
}

// InsertTransaction mocks base method.
func (m *MockApp) InsertTransaction(ctx context.Context, transaction *model.Transaction) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransaction", ctx, transaction)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTransaction indicates an expected call of InsertTransaction.
func (mr *MockAppMockRecorder) InsertTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransaction", reflect.TypeOf((*MockApp)(nil).InsertTransaction), ctx, transaction)
}

// InsertTransactions mocks base method.
func (m *MockApp) InsertTransactions(ctx context.Context, transaction []*model.Transaction) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransactions", ctx, transaction)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTransactions indicates an expected call of InsertTransactions.
func (mr *MockAppMockRecorder) InsertTransactions(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransactions", reflect.TypeOf((*MockApp)(nil).InsertTransactions), ctx, transaction)
}
