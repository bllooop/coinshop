package repository

import (
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	SignUser(username, password string) (domain.User, error)
}
type Shop interface {
	BuyItem(userid int, name string) (int, error)
	SendCoin(userid int, input domain.Transactions) (int, error)
}

type Repository struct {
	Authorization
	Shop
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Shop:          NewShopPostgres(db),
	}
}
