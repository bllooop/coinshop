package usecase

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
	"github.com/golang-jwt/jwt"
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
	return s.repo.CreateUser(user)
}

func (s *AuthUsecase) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.SignUser(username, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
	/*
			token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user.Id
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	*/
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

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
