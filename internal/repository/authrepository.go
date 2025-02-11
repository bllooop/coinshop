package repository

import (
	"context"
	"fmt"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	pg *pgxpool.Pool
}

func NewAuthPostgres(pg *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{
		pg: pg,
	}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (username,password,coins) VALUES ($1,$2,$3) RETURNING id`, userListTable)
	row := r.pg.QueryRow(context.Background(), query, user.UserName, user.Password, user.Coins)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) SignUser(username, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT id,username,password FROM %s WHERE username=$1 AND password=$2`, userListTable)
	res := r.pg.QueryRow(context.Background(), query, username, password)
	err := res.Scan(&user.Id, &user.UserName, &user.Password)

	return user, err
}
