package usersvc_test

import (
	"context"
	"errors"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/core/ports/mocks"
	"robinhood/internal/core/services/usersvc"
	"robinhood/internal/errmsg"
	"robinhood/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testModule struct {
	ur  *mocks.UserRepository
	svc ports.UserService
}

type test struct {
	name     string
	args     []interface{}
	mockFn   func(*testModule)
	assertFn func(*testModule)
}

var (
	ctx = context.TODO()
)

func new(t *testing.T) *testModule {
	ur := mocks.NewUserRepository(t)
	return &testModule{
		ur:  ur,
		svc: usersvc.New(ur),
	}
}

func TestRegister(t *testing.T) {
	var err error
	mockReq := &domains.RegisterRequest{
		Username: "username",
		Password: "password",
	}

	tests := []*test{
		{
			name: "return error when get user by username failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(nil, errors.New("error"))
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.Equal(t, errmsg.UserRegisterFailed, err)
			},
		},
		{
			name: "return error when user already exists",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(&domains.User{}, nil)
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.Equal(t, errmsg.UserExisted, err)
			},
		},
		{
			name: "return error when create user failed",
			args: []interface{}{
				ctx,
				&domains.RegisterRequest{
					Username: "username",
					Password: "password",
					Email:    "email",
				},
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(nil, nil)
				m.ur.On("Create", ctx, mock.AnythingOfType("*domains.CreateUserRequest")).Return(nil, errors.New("error"))
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.Equal(t, errmsg.UserRegisterFailed, err)
			},
		},
		{
			name: "success",
			args: []interface{}{
				ctx,
				&domains.RegisterRequest{
					Username: "username",
					Password: "password",
					Email:    "email",
				},
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(nil, nil)
				m.ur.On("Create", ctx, mock.AnythingOfType("*domains.CreateUserRequest")).Return(&domains.User{}, nil)
			},
			assertFn: func(m *testModule) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(t)
			tc.mockFn(m)
			err = m.svc.Register(tc.args[0].(context.Context), tc.args[1].(*domains.RegisterRequest))
			tc.assertFn(m)
		})
	}
}

func TestLogin(t *testing.T) {
	var result *domains.LoginResponse
	var err error
	mockReq := &domains.LoginRequest{
		Username: "username",
		Password: "password",
	}

	tests := []*test{
		{
			name: "return error when get user by username failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(nil, errors.New("error"))
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.Equal(t, errmsg.UserLoginFailed, err)
			},
		},
		{
			name: "return error when user not found",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(nil, nil)
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.Equal(t, errmsg.UsernameOrPasswordIncorrect, err)
			},
		},
		{
			name: "should error when password incorrect",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(&domains.User{
					Username: mockReq.Username,
					Password: mockReq.Password,
				}, nil)
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.Equal(t, errmsg.UsernameOrPasswordIncorrect, err)
			},
		},
		{
			name: "success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				hash, _ := utils.HashPassword(mockReq.Password, utils.DefaultCost)
				m.ur.On("GetByUsername", ctx, mockReq.Username).Return(&domains.User{
					Username: mockReq.Username,
					Password: hash,
				}, nil)
			},
			assertFn: func(m *testModule) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(t)
			tc.mockFn(m)
			result, err = m.svc.Login(tc.args[0].(context.Context), tc.args[1].(*domains.LoginRequest))
			tc.assertFn(m)
		})
	}
}

func TestUpdate(t *testing.T) {
	var result *domains.User
	var err error
	mockReq := &domains.UpdateUserRequest{
		UserId:       "user_id",
		ProfileImage: "profile_image",
	}

	tests := []*test{
		{
			name: "return error when update failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("Update", ctx, mockReq).Return(nil, errors.New("error"))
			},
			assertFn: func(m *testModule) {
				assert.Error(t, err)
				assert.EqualError(t, err, "error")
			},
		},
		{
			name: "success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(m *testModule) {
				m.ur.On("Update", ctx, mockReq).Return(&domains.User{}, nil)
			},
			assertFn: func(m *testModule) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(t)
			tc.mockFn(m)
			result, err = m.svc.Update(tc.args[0].(context.Context), tc.args[1].(*domains.UpdateUserRequest))
			tc.assertFn(m)
		})
	}
}
