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

func TestHandler_sendCoin(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockShop, userid int, transactions domain.Transactions)

	testTable := []struct {
		name                 string
		inputBody            string
		inputTransactions    domain.Transactions
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputBody:   `{"destination":2, "amount":10}`,
			inputUserId: 1,
			inputTransactions: domain.Transactions{
				Destination: 2,
				Amount:      10,
			},
			mockBehavior: func(s *mock_usecase.MockShop, userid int, transactions domain.Transactions) {
				s.EXPECT().SendCoin(userid, transactions).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:        "Error during execution in service",
			inputBody:   `{"destination":2, "amount":10}`,
			inputUserId: 1,
			inputTransactions: domain.Transactions{
				Destination: 2,
				Amount:      10,
			},
			mockBehavior: func(s *mock_usecase.MockShop, userid int, transactions domain.Transactions) {
				s.EXPECT().SendCoin(userid, transactions).Return(0, errors.New("Internal Server Error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Internal Server Error"}`,
		},
		{
			name:        "Bad input",
			inputBody:   `{"amount":-100, "destination":1}`,
			inputUserId: 1,
			inputTransactions: domain.Transactions{
				Amount:      -100,
				Destination: 1,
			},
			mockBehavior:         func(s *mock_usecase.MockShop, userid int, transactions domain.Transactions) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Значения получателя и суммы не могут быть отрицательными"}`,
		},
		{
			name:        "Missing destination",
			inputBody:   `{"amount":10}`,
			inputUserId: 1,
			inputTransactions: domain.Transactions{
				Amount: 10,
			},
			mockBehavior:         func(s *mock_usecase.MockShop, userid int, transactions domain.Transactions) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'Transactions.Destination' Error:Field validation for 'Destination' failed on the 'required' tag"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockShop(c)
			testCase.mockBehavior(repo, testCase.inputUserId, testCase.inputTransactions)

			usecases := &usecase.Usecase{Shop: repo}
			handler := Handler{usecases}
			r := gin.New()
			r.POST("/api/sendCoin", func(c *gin.Context) {
				c.Set("userId", testCase.inputUserId)
				handler.sendCoin(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/sendCoin",
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

func TestHandler_buyItem(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockShop, name string, userId int)

	testTable := []struct {
		name                 string
		inputName            string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputName:   "cup",
			inputUserId: 1,
			mockBehavior: func(s *mock_usecase.MockShop, name string, userId int) {
				s.EXPECT().BuyItem(userId, name).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:        "Error during execution in service",
			inputName:   "cup",
			inputUserId: 1,
			mockBehavior: func(s *mock_usecase.MockShop, name string, userId int) {
				s.EXPECT().BuyItem(userId, name).Return(0, errors.New("Internal Server Error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Internal Server Error"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockShop(c)
			testCase.mockBehavior(repo, testCase.inputName, testCase.inputUserId)

			usecases := &usecase.Usecase{Shop: repo}
			handler := Handler{usecases}
			r := gin.New()
			r.PUT("/api/buy/:name", func(c *gin.Context) {
				c.Set("userId", testCase.inputUserId)
				handler.buyItem(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/buy/"+testCase.inputName, nil)

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

func intPointer(s int) *int {
	return &s
}
