// Code generated by MockGen. DO NOT EDIT.
// Source: ./database/db_types.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	database "github.com/chapdast/project_chat/database"
	gomock "github.com/golang/mock/gomock"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// MockProjectsDB is a mock of ProjectsDB interface.
type MockProjectsDB struct {
	ctrl     *gomock.Controller
	recorder *MockProjectsDBMockRecorder
}

// MockProjectsDBMockRecorder is the mock recorder for MockProjectsDB.
type MockProjectsDBMockRecorder struct {
	mock *MockProjectsDB
}

// NewMockProjectsDB creates a new mock instance.
func NewMockProjectsDB(ctrl *gomock.Controller) *MockProjectsDB {
	mock := &MockProjectsDB{ctrl: ctrl}
	mock.recorder = &MockProjectsDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectsDB) EXPECT() *MockProjectsDBMockRecorder {
	return m.recorder
}

// HaveAccess mocks base method.
func (m *MockProjectsDB) HaveAccess(ctx context.Context, userId, projectId uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HaveAccess", ctx, userId, projectId)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HaveAccess indicates an expected call of HaveAccess.
func (mr *MockProjectsDBMockRecorder) HaveAccess(ctx, userId, projectId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HaveAccess", reflect.TypeOf((*MockProjectsDB)(nil).HaveAccess), ctx, userId, projectId)
}

// MockProjectChatDB is a mock of ProjectChatDB interface.
type MockProjectChatDB struct {
	ctrl     *gomock.Controller
	recorder *MockProjectChatDBMockRecorder
}

// MockProjectChatDBMockRecorder is the mock recorder for MockProjectChatDB.
type MockProjectChatDBMockRecorder struct {
	mock *MockProjectChatDB
}

// NewMockProjectChatDB creates a new mock instance.
func NewMockProjectChatDB(ctrl *gomock.Controller) *MockProjectChatDB {
	mock := &MockProjectChatDB{ctrl: ctrl}
	mock.recorder = &MockProjectChatDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectChatDB) EXPECT() *MockProjectChatDBMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockProjectChatDB) Read(ctx context.Context, projectId uint64) ([]*database.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, projectId)
	ret0, _ := ret[0].([]*database.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockProjectChatDBMockRecorder) Read(ctx, projectId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockProjectChatDB)(nil).Read), ctx, projectId)
}

// Save mocks base method.
func (m *MockProjectChatDB) Save(ctx context.Context, message *database.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockProjectChatDBMockRecorder) Save(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockProjectChatDB)(nil).Save), ctx, message)
}

// MockProjectManager is a mock of ProjectManager interface.
type MockProjectManager struct {
	ctrl     *gomock.Controller
	recorder *MockProjectManagerMockRecorder
}

// MockProjectManagerMockRecorder is the mock recorder for MockProjectManager.
type MockProjectManagerMockRecorder struct {
	mock *MockProjectManager
}

// NewMockProjectManager creates a new mock instance.
func NewMockProjectManager(ctrl *gomock.Controller) *MockProjectManager {
	mock := &MockProjectManager{ctrl: ctrl}
	mock.recorder = &MockProjectManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectManager) EXPECT() *MockProjectManagerMockRecorder {
	return m.recorder
}

// Client mocks base method.
func (m *MockProjectManager) Client() *mongo.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Client")
	ret0, _ := ret[0].(*mongo.Client)
	return ret0
}

// Client indicates an expected call of Client.
func (mr *MockProjectManagerMockRecorder) Client() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Client", reflect.TypeOf((*MockProjectManager)(nil).Client))
}

// HaveAccess mocks base method.
func (m *MockProjectManager) HaveAccess(ctx context.Context, userId, projectId uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HaveAccess", ctx, userId, projectId)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HaveAccess indicates an expected call of HaveAccess.
func (mr *MockProjectManagerMockRecorder) HaveAccess(ctx, userId, projectId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HaveAccess", reflect.TypeOf((*MockProjectManager)(nil).HaveAccess), ctx, userId, projectId)
}

// Read mocks base method.
func (m *MockProjectManager) Read(ctx context.Context, projectId uint64) ([]*database.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, projectId)
	ret0, _ := ret[0].([]*database.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockProjectManagerMockRecorder) Read(ctx, projectId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockProjectManager)(nil).Read), ctx, projectId)
}

// Save mocks base method.
func (m *MockProjectManager) Save(ctx context.Context, message *database.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockProjectManagerMockRecorder) Save(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockProjectManager)(nil).Save), ctx, message)
}