package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bllooop/coinshop/internal/delivery/api"
	"github.com/bllooop/coinshop/internal/repository"
	"github.com/bllooop/coinshop/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *PostgresContainer
	db          *sqlx.DB
	repository  *repository.Repository
	handler     *api.Handler
}

func (suite *AuthHandlerTestSuite) SetupSuite() {
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

	usecases := &usecase.Usecase{
		Authorization: usecase.NewAuthUsecase(suite.repository),
	}

	suite.handler = &api.Handler{Usecases: usecases}
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	_, err := suite.db.Exec("TRUNCATE TABLE userlist RESTART IDENTITY CASCADE")
	assert.NoError(suite.T(), err)
}

func (suite *AuthHandlerTestSuite) TearDownSuite() {
	assert.NoError(suite.T(), suite.pgContainer.Terminate(suite.ctx))
}

func (suite *AuthHandlerTestSuite) TestSignUp() {

	reqBody := `{"username":"name", "password":"password1"}`
	r := gin.New()
	r.POST("/api/auth/sign-up", suite.handler.SignUp)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/sign-up",
		bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	expectedResponse := `{"id":1}`
	assert.JSONEq(suite.T(), expectedResponse, w.Body.String())

	var user string
	err := suite.db.QueryRow("SELECT username FROM userlist WHERE id = 1").Scan(&user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "name", user)
}

func (suite *AuthHandlerTestSuite) TestSignIn() {
	pass, err := usecase.HashPassword("password1")
	assert.NoError(suite.T(), err)
	_, err = suite.db.Exec("INSERT INTO userlist (username, coins, password) VALUES ($1, $2, $3)",
		"name", 1000, pass)
	assert.NoError(suite.T(), err)
	reqBody := `{"username":"name", "password":"password1"}`
	r := gin.New()
	r.POST("/api/auth/sign-in", suite.handler.SignIn)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/sign-in",
		bytes.NewBufferString(reqBody))
	r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	token, exists := response["token"].(string)
	assert.True(suite.T(), exists, "Response should contain a token")
	assert.NotEmpty(suite.T(), token, "Token should not be empty")

	parts := strings.Split(token, ".")
	assert.Equal(suite.T(), 3, len(parts), "JWT should have 3 parts")

	var username string
	err = suite.db.QueryRow("SELECT username FROM userlist WHERE id = 1").Scan(&username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "name", username)
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
