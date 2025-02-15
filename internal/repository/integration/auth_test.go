package integration

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthRepoTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *PostgresContainer
	repository  *repository.AuthPostgres
	db          *sqlx.DB
}

func (suite *AuthRepoTestSuite) SetupSuite() {
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
	suite.repository = repository.NewAuthPostgres(suite.db)

}
func (suite *AuthRepoTestSuite) SetupTest() {
	_, err := suite.repository.DB().Exec("TRUNCATE TABLE userlist RESTART IDENTITY CASCADE")
	assert.NoError(suite.T(), err)
}
func (suite *AuthRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("error terminating postgres container")
	}
}

func (suite *AuthRepoTestSuite) TestCreateUser() {
	t := suite.T()
	input := domain.User{
		UserName: "name1",
		Password: "password123",
		Coins:    IntPointer(1000),
	}
	//_, err := suite.repository.DB().Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
	//	"name", 1000, "password123")
	//assert.NoError(t, err)

	createUser, err := suite.repository.CreateUser(input)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}
	assert.NotNil(t, createUser)
	assert.Greater(t, createUser, 0)

	var inputCheck domain.User
	err = suite.repository.DB().Get(&inputCheck, "SELECT id, username, password FROM userlist WHERE id = $1", createUser)
	assert.NoError(t, err)
	assert.Equal(t, IntPointer(1), inputCheck.Id)
	assert.Equal(t, input.UserName, inputCheck.UserName)
	assert.Equal(t, input.Password, inputCheck.Password)

}

func (suite *AuthRepoTestSuite) TestGetUser() {
	t := suite.T()
	inputUsername := "name"
	_, err := suite.repository.DB().Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		"name", 1000, "password123")
	assert.NoError(t, err)

	getUser, err := suite.repository.SignUser(inputUsername)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}
	assert.NotNil(t, getUser)

	var inputCheck domain.User
	err = suite.repository.DB().Get(&inputCheck, "SELECT id, username, password FROM userlist WHERE id = $1", getUser)
	assert.NoError(t, err)
	assert.Equal(t, getUser.Id, inputCheck.Id)
	assert.Equal(t, getUser.UserName, inputCheck.UserName)
	assert.Equal(t, getUser.Password, inputCheck.Password)

}
