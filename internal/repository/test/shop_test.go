package repository

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShopRepoTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *PostgresContainer
	repository  *repository.ShopPostgres
	db          *sqlx.DB
}

func (suite *ShopRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := CreatePostgresContainer(suite.ctx)
	if err != nil {
		logger.Log.Fatal().Err(err)
	}
	host, err := pgContainer.Host(suite.ctx)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to get container host")
	}
	port, err := pgContainer.MappedPort(suite.ctx, "5432")
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to get mapped port")
	}

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
	if err != nil {
		logger.Log.Fatal().Err(err)
	}
	suite.db = db
	migratePath, err := filepath.Abs("../../../migrations")
	fmt.Println("DEBUG: Running migrations from", migratePath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Migration failed")
	}
	err = repository.RunMigrate(cfg, migratePath)

	if err != nil {
		logger.Log.Fatal().Err(err)

	}
	suite.repository = repository.NewShopPostgres(suite.db)

}

func (suite *ShopRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("error terminating postgres container")
	}
}

func (suite *ShopRepoTestSuite) TestSendingCoin() {
	t := suite.T()
	_, err := suite.repository.DB().Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		"name", 1000, "password123")
	assert.NoError(t, err)
	_, err = suite.repository.DB().Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		"name2", 1000, "password123")
	assert.NoError(t, err)
	input := domain.Transactions{
		Source:              IntPointer(2),
		DestinationUsername: "name",
		Destination:         IntPointer(1),
		Amount:              10,
		Timestamp:           func() *time.Time { t := time.Now(); return &t }(),
	}
	sendCoin, err := suite.repository.SendCoin(input)
	if err != nil {
		t.Fatal("Failed to send coin:", err)
	}
	assert.NotNil(t, sendCoin)
	assert.Greater(t, sendCoin, 0)

	var transactionCount int
	err = suite.repository.DB().Get(&transactionCount, "SELECT COUNT(*) FROM transactions WHERE id = $1", sendCoin)
	if err != nil {
		t.Fatal("Error checking transaction count:", err)
	}
	assert.Equal(t, 1, transactionCount)

}

func (suite *ShopRepoTestSuite) TestBuyingItem() {
	t := suite.T()
	_, err := suite.repository.DB().Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		"name", 1000, "password123")
	assert.NoError(t, err)
	_, err = suite.repository.DB().Exec("INSERT INTO shop (name, price) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		"cup", 20)
	assert.NoError(t, err)
	inputUserid := 1
	inputName := "cup"
	sendCoin, err := suite.repository.BuyItem(inputUserid, inputName)
	if err != nil {
		t.Fatal("Failed to send coin:", err)
	}
	assert.NotNil(t, sendCoin)
	assert.Greater(t, sendCoin, 0)

	var purchasesCount int
	err = suite.repository.DB().Get(&purchasesCount, "SELECT COUNT(*) FROM purchases WHERE id = $1", sendCoin)
	if err != nil {
		t.Fatal("Error checking purchases count:", err)
	}
	assert.Equal(t, 1, purchasesCount)

}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ShopRepoTestSuite))
}

/*func TestInsertAndRetrieveUser(t *testing.T) {
	_, err := testDB.Exec(`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = testDB.Exec(`INSERT INTO users (name) VALUES ($1)`, "Alice")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	var name string
	err = testDB.QueryRow(`SELECT name FROM users WHERE name = $1`, "Alice").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", name)
	}
}
*/

func IntPointer(s int) *int {
	return &s
}
