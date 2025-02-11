package usecase

import (
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}
type Shop interface {
	BuyItem(userid int, name string) (int, error)
	SendCoin(input domain.Transactions) (int, error)
}
type Usecase struct {
	Authorization
	Shop
}

func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{
		Authorization: NewAuthUsecase(repo),
		Shop:          NewShopUsecase(repo),
	}
}
