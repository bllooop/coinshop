package repository

import (
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	SignUser(username, password string) (domain.User, error)
}
type Shop interface {
	BuyItem(userid int, name string) (int, error)
	SendCoin(input domain.Transactions) (int, error)
}

type Repository struct {
	Authorization
	Shop
}

func NewRepository(pg *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(pg),
		Shop:          NewShopPostgres(pg),
	}
}
