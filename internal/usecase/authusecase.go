package usecase

import (
	"errors"
	"time"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo repository.Authorization
}

func NewAuthUsecase(repo *repository.Repository) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *AuthUsecase) CreateUser(user domain.User) (int, error) {
	var err error
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateUser(user)
}
func (s *AuthUsecase) SignUser(username, password string) (domain.User, error) {
	user, err := s.repo.SignUser(username)
	if err != nil {
		return domain.User{}, err
	}
	if !verifyPassword(user.Password, password) {
		return domain.User{}, errors.New("invalid credentials")
	}
	return user, nil
}
func (s *AuthUsecase) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthUsecase) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("некорретный signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims не типа *tokenClaims")
	}

	return claims.UserId, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
