package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/usecase"
	mock_usecase "github.com/bllooop/coinshop/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockAuthorization, user domain.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"12345"}`,
			inputUser: domain.User{
				UserName: "test",
				Password: "12345",
				Coins:    intPointer(1000),
			},
			mockBehavior: func(s *mock_usecase.MockAuthorization, user domain.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Error during execution in service",
			inputBody: `{"username": "test", "password":"12345"}`,
			inputUser: domain.User{
				UserName: "test",
				Password: "12345",
				Coins:    intPointer(1000),
			},
			mockBehavior: func(s *mock_usecase.MockAuthorization, user domain.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("Internal Server Error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Internal Server Error"}`,
		},
		{
			name:                 "Bad input",
			inputBody:            `{"username":1000}`,
			inputUser:            domain.User{},
			mockBehavior:         func(s *mock_usecase.MockAuthorization, user domain.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"json: cannot unmarshal number into Go struct field User.username of type string"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockAuthorization(c)
			testCase.mockBehavior(repo, testCase.inputUser)

			usecases := &usecase.Usecase{Authorization: repo}
			handler := Handler{usecases}
			r := gin.New()
			api := r.Group("/api")
			auth := api.Group("/auth")
			auth.POST("/sign-up", handler.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/auth/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			if json.Valid([]byte(testCase.expectedResponseBody)) {
				assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
			} else {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockAuthorization, username, password string)

	testTable := []struct {
		name                 string
		inputBody            string
		username             string
		password             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful SignIn",
			inputBody: `{"username":"name", "password":"12345"}`,
			username:  "test",
			password:  "12345",
			mockBehavior: func(s *mock_usecase.MockAuthorization, username, password string) {
				s.EXPECT().SignUser(username, password).Return(domain.User{Id: 1}, nil)
				s.EXPECT().GenerateToken(1).Return("valid.jwt.token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"valid.jwt.token"}`,
		},
		{
			name:      "User Not Found - New Account Created",
			inputBody: `{"username":"notname", "password":"password123"}`,
			username:  "newuser",
			password:  "password123",
			mockBehavior: func(s *mock_usecase.MockAuthorization, username, password string) {
				s.EXPECT().SignUser(username, password).Return(domain.User{}, sql.ErrNoRows)
				s.EXPECT().CreateUser(domain.User{UserName: username, Password: password, Coins: intPointer(1000)}).Return(2, nil)
				s.EXPECT().GenerateToken(2).Return("newuser.jwt.token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":2, "token":"newuser.jwt.token"}`,
		},
		{
			name:                 "Invalid JSON Input",
			inputBody:            `{"name":1000}`,
			mockBehavior:         func(s *mock_usecase.MockAuthorization, username, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"json: cannot unmarshal number into Go struct field SignInInput.username of type string"}`,
		},
		{
			name:      "SignUser Error",
			inputBody: `{"name":"test", "password":"12345"}`,
			username:  "test",
			password:  "12345",
			mockBehavior: func(s *mock_usecase.MockAuthorization, username, password string) {
				s.EXPECT().SignUser(username, password).Return(domain.User{}, errors.New("Internal Server Error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Error checking user: Internal Server Error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockAuthorization(c)
			testCase.mockBehavior(repo, testCase.username, testCase.password)

			usecases := &usecase.Usecase{Authorization: repo}
			handler := Handler{usecases}
			r := gin.New()
			r.POST("/api/auth/sign-in", handler.SignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/auth/sign-in",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			if json.Valid([]byte(testCase.expectedResponseBody)) {
				assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
			} else {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
