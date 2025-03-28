// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/mock.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	domain "github.com/bllooop/coinshop/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
	isgomock struct{}
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(user domain.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), user)
}

// GenerateToken mocks base method.
func (m *MockAuthorization) GenerateToken(userId int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthorizationMockRecorder) GenerateToken(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthorization)(nil).GenerateToken), userId)
}

// ParseToken mocks base method.
func (m *MockAuthorization) ParseToken(accessToken string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", accessToken)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAuthorizationMockRecorder) ParseToken(accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAuthorization)(nil).ParseToken), accessToken)
}

// SignUser mocks base method.
func (m *MockAuthorization) SignUser(username, password string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUser", username, password)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUser indicates an expected call of SignUser.
func (mr *MockAuthorizationMockRecorder) SignUser(username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUser", reflect.TypeOf((*MockAuthorization)(nil).SignUser), username, password)
}

// MockShop is a mock of Shop interface.
type MockShop struct {
	ctrl     *gomock.Controller
	recorder *MockShopMockRecorder
	isgomock struct{}
}

// MockShopMockRecorder is the mock recorder for MockShop.
type MockShopMockRecorder struct {
	mock *MockShop
}

// NewMockShop creates a new mock instance.
func NewMockShop(ctrl *gomock.Controller) *MockShop {
	mock := &MockShop{ctrl: ctrl}
	mock.recorder = &MockShopMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShop) EXPECT() *MockShopMockRecorder {
	return m.recorder
}

// BuyItem mocks base method.
func (m *MockShop) BuyItem(userid int, name string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyItem", userid, name)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuyItem indicates an expected call of BuyItem.
func (mr *MockShopMockRecorder) BuyItem(userid, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyItem", reflect.TypeOf((*MockShop)(nil).BuyItem), userid, name)
}

// GetUserSummary mocks base method.
func (m *MockShop) GetUserSummary(userID int) (*domain.UserSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSummary", userID)
	ret0, _ := ret[0].(*domain.UserSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSummary indicates an expected call of GetUserSummary.
func (mr *MockShopMockRecorder) GetUserSummary(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSummary", reflect.TypeOf((*MockShop)(nil).GetUserSummary), userID)
}

// SendCoin mocks base method.
func (m *MockShop) SendCoin(userid int, input domain.Transactions) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoin", userid, input)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendCoin indicates an expected call of SendCoin.
func (mr *MockShopMockRecorder) SendCoin(userid, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoin", reflect.TypeOf((*MockShop)(nil).SendCoin), userid, input)
}
