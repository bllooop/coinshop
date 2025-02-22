package usecase

import (
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mock.go
type Authorization interface {
	CreateUser(user domain.User) (int, error)
	SignUser(username, password string) (domain.User, error)
	GenerateToken(userId int) (string, error)
	ParseToken(accessToken string) (int, error)
}
type Shop interface {
	BuyItem(userid int, name string) (int, error)
	SendCoin(userid int, input domain.Transactions) (int, error)
	GetUserSummary(userID int) (*domain.UserSummary, error)
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
