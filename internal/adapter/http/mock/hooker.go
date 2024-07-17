// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go
//
// Generated by this command:
//
//	mockgen -source=contract.go -destination=mock/hooker.go -package=mock hookService
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hooker "github.com/k1nky/tookhook/internal/service/hooker"
	gomock "go.uber.org/mock/gomock"
)

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Debugf mocks base method.
func (m *Mocklogger) Debugf(template string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockloggerMockRecorder) Debugf(template any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*Mocklogger)(nil).Debugf), varargs...)
}

// Errorf mocks base method.
func (m *Mocklogger) Errorf(template string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockloggerMockRecorder) Errorf(template any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*Mocklogger)(nil).Errorf), varargs...)
}

// Infof mocks base method.
func (m *Mocklogger) Infof(template string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockloggerMockRecorder) Infof(template any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*Mocklogger)(nil).Infof), varargs...)
}

// MockhookService is a mock of hookService interface.
type MockhookService struct {
	ctrl     *gomock.Controller
	recorder *MockhookServiceMockRecorder
}

// MockhookServiceMockRecorder is the mock recorder for MockhookService.
type MockhookServiceMockRecorder struct {
	mock *MockhookService
}

// NewMockhookService creates a new mock instance.
func NewMockhookService(ctrl *gomock.Controller) *MockhookService {
	mock := &MockhookService{ctrl: ctrl}
	mock.recorder = &MockhookServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockhookService) EXPECT() *MockhookServiceMockRecorder {
	return m.recorder
}

// Forward mocks base method.
func (m *MockhookService) Forward(ctx context.Context, name string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Forward", ctx, name, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Forward indicates an expected call of Forward.
func (mr *MockhookServiceMockRecorder) Forward(ctx, name, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Forward", reflect.TypeOf((*MockhookService)(nil).Forward), ctx, name, data)
}

// Reload mocks base method.
func (m *MockhookService) Reload(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reload", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reload indicates an expected call of Reload.
func (mr *MockhookServiceMockRecorder) Reload(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*MockhookService)(nil).Reload), ctx)
}

// Status mocks base method.
func (m *MockhookService) Status(ctx context.Context) hooker.ServiceStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", ctx)
	ret0, _ := ret[0].(hooker.ServiceStatus)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockhookServiceMockRecorder) Status(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockhookService)(nil).Status), ctx)
}
