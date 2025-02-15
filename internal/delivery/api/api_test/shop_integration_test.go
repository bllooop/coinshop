package api_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/bllooop/coinshop/internal/delivery/api"
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
	"github.com/bllooop/coinshop/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SendCoinHandlerTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *PostgresContainer
	db          *sqlx.DB
	repository  *repository.Repository
	handler     *api.Handler
}

func (suite *SendCoinHandlerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := CreatePostgresContainer(suite.ctx)
	assert.NoError(suite.T(), err)

	host, err := pgContainer.Host(suite.ctx)
	assert.NoError(suite.T(), err)
	port, err := pgContainer.MappedPort(suite.ctx, "5432")
	assert.NoError(suite.T(), err)

	cfg := repository.Config{
		Username: "postgres",
		Password: "postgres",
		Host:     host,
		Port:     port.Port(),
		DBname:   "test-db",
		SSLMode:  "disable",
	}

	suite.pgContainer = pgContainer
	db, err := repository.NewPostgresDB(cfg)
	assert.NoError(suite.T(), err)
	suite.db = db

	migratePath, err := filepath.Abs("../../../../migrations")
	assert.NoError(suite.T(), err)
	fmt.Println("DEBUG: Running migrations from", migratePath)
	err = repository.RunMigrate(cfg, migratePath)
	assert.NoError(suite.T(), err)

	suite.repository = repository.NewRepository(db)

	// Initialize the usecase with the repository
	usecases := &usecase.Usecase{
		Shop: usecase.NewShopUsecase(suite.repository),
	}

	suite.handler = &api.Handler{Usecases: usecases}
}

func (suite *SendCoinHandlerTestSuite) SetupTest() {
	_, err := suite.db.Exec("TRUNCATE TABLE userlist, transactions,purchases RESTART IDENTITY CASCADE")
	assert.NoError(suite.T(), err)
}

func (suite *SendCoinHandlerTestSuite) TearDownSuite() {
	assert.NoError(suite.T(), suite.pgContainer.Terminate(suite.ctx))
}

func (suite *SendCoinHandlerTestSuite) TestSendCoin() {

	_, err := suite.db.Exec("INSERT INTO userlist (id, username, password, coins) VALUES (1, 'name1','password123', 1000), (2, 'name2','password123', 0)")
	assert.NoError(suite.T(), err)

	reqBody := `{"destination_username":"name2", "amount":10}`
	r := gin.New()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("userId", 1)
		suite.handler.SendCoin(c)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/sendCoin",
		bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	expectedResponse := `{"id":1}`
	assert.JSONEq(suite.T(), expectedResponse, w.Body.String())

	var transaction domain.Transactions
	err = suite.db.QueryRow("SELECT id, amount, source, destination FROM transactions WHERE id = 1").
		Scan(&transaction.Id, &transaction.Amount, &transaction.Source, &transaction.Destination)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, transaction.Amount)
	assert.Equal(suite.T(), IntPointer(1), transaction.Source)
	assert.Equal(suite.T(), IntPointer(2), transaction.Destination)

	var senderCoins int
	err = suite.db.QueryRow("SELECT coins FROM userlist WHERE id = 1").Scan(&senderCoins)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 990, senderCoins)

	var recipientCoins int
	err = suite.db.QueryRow("SELECT coins FROM userlist WHERE id = 2").Scan(&recipientCoins)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, recipientCoins)
}

func (suite *SendCoinHandlerTestSuite) TestBuyItem() {

	_, err := suite.db.Exec("INSERT INTO userlist (id, username, password, coins) VALUES (1, 'name1','password123', 1000)")
	assert.NoError(suite.T(), err)
	inputName := "t-shirt"
	r := gin.New()
	r.PUT("/api/buy/:item", func(c *gin.Context) {
		c.Set("userId", 1)
		suite.handler.BuyItem(c)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/api/buy/"+inputName,
		nil)
	r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	expectedResponse := `{"id":1}`
	assert.JSONEq(suite.T(), expectedResponse, w.Body.String())

	var userID, itemID int
	err = suite.db.QueryRow("SELECT user_id, item_id FROM purchases WHERE id = 1").Scan(&userID, &itemID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, userID)
	assert.Equal(suite.T(), 1, itemID)
	var buyerCoins int
	err = suite.db.QueryRow("SELECT coins FROM userlist WHERE id = 1").Scan(&buyerCoins)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 920, buyerCoins)
}

func (suite *SendCoinHandlerTestSuite) TestGetInfo() {

	_, err := suite.db.Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		"name", 1000, "password123")
	assert.NoError(suite.T(), err)
	_, err = suite.db.Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		"name2", 1000, "password123")
	assert.NoError(suite.T(), err)
	_, err = suite.db.Exec("INSERT INTO transactions (source, destination,amount,transaction_time) VALUES ($1, $2,$3,$4) ON CONFLICT (id) DO NOTHING",
		1, 2, 10, time.Now())
	assert.NoError(suite.T(), err)
	_, err = suite.db.Exec("INSERT INTO shop (name, price) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		"t-shirt", 80)
	assert.NoError(suite.T(), err)
	_, err = suite.db.Exec("INSERT INTO purchases (user_id, item_id, price,purchase_date) VALUES ($1, $2,$3,$4) ON CONFLICT (id) DO NOTHING",
		1, 1, 80, time.Now())
	assert.NoError(suite.T(), err)
	_, err = suite.db.Exec("UPDATE userlist SET coins = coins - 90 WHERE username = $1", "name")
	assert.NoError(suite.T(), err)
	r := gin.New()
	r.GET("/api/info", func(c *gin.Context) {
		c.Set("userId", 1)
		suite.handler.GetInfo(c)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/info", nil)
	r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	expectedResponse := `{
					"username": "name",
					"coins": 910,
					"purchased_items": [{ "item_name": "t-shirt", "quantity": 1 }],
					"transactions_summary": {
						"received_coins": null,
						"sent_coins": [
							{
								"destination": 2,
								"destination_username": "name2",
								"amount": 10
							}
						]
					}
				}`
	assert.JSONEq(suite.T(), expectedResponse, w.Body.String())

}

func TestSendCoinHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SendCoinHandlerTestSuite))
}

func IntPointer(s int) *int {
	return &s
}
