package api

import (
	"bytes"
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
			name:      "Ошибка выполнения запроса",
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
			name:                 "Плохой ввод",
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
			name:      "OK",
			inputBody: `{"username":"name", "password":"12345"}`,
			username:  "name",
			password:  "12345",
			mockBehavior: func(s *mock_usecase.MockAuthorization, username, password string) {
				s.EXPECT().SignUser("name", "12345").Return(domain.User{Id: 1}, nil)
				s.EXPECT().GenerateToken(1).Return("valid.jwt.token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"valid.jwt.token"}`,
		},
		{
			name:      "Пользователь не найден - регистрация",
			inputBody: `{"username":"notname", "password":"password123"}`,
			username:  "notname",
			password:  "password123",
			mockBehavior: func(s *mock_usecase.MockAuthorization, username, password string) {
				s.EXPECT().SignUser("notname", "password123").Return(domain.User{}, errors.New("пользователь не найден"))
				s.EXPECT().CreateUser(domain.User{UserName: "notname", Password: "password123", Coins: intPointer(1000)}).Return(2, nil)
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
			expectedResponseBody: `{"message":"Key: 'SignInInput.UserName' Error:Field validation for 'UserName' failed on the 'required' tag\nKey: 'SignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Ошибка авторизации",
			inputBody: `{"username":"test", "password":"12345"}`,
			username:  "test",
			password:  "12345",
			mockBehavior: func(s *mock_usecase.MockAuthorization, username, password string) {
				s.EXPECT().SignUser("test", "12345").Return(domain.User{}, errors.New("Internal Server Error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Ошибка авторизации: Internal Server Error"}`,
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
			req := httptest.NewRequest("POST", "/api/auth/sign-in", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			// Check if the response body is valid JSON
			if json.Valid([]byte(testCase.expectedResponseBody)) {
				assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
			} else {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
