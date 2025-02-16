package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (username,password,coins) VALUES ($1,$2,$3) RETURNING id`, userListTable)
	row := r.db.QueryRowx(query, user.UserName, user.Password, user.Coins)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) SignUser(username string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT id,username,password FROM %s WHERE username=$1`, userListTable)
	res := r.db.QueryRowx(query, username)
	err := res.Scan(&user.Id, &user.UserName, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("пользователь не найден")
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *AuthPostgres) DB() *sqlx.DB {
	return r.db
}
