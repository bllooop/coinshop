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
			auth.POST("/sign-up", handler.signUp)

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
