package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/bllooop/coinshop/internal/domain"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type ShopPostgres struct {
	db *sqlx.DB
}

func NewShopPostgres(db *sqlx.DB) *ShopPostgres {
	return &ShopPostgres{
		db: db,
	}
}

func (r *ShopPostgres) BuyItem(userid int, name string) (int, error) {
	tr, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tr.Rollback()
	var id, itemID, price, amount int
	getIdQuery := fmt.Sprintf("SELECT id, price FROM %s WHERE name = $1", shopTable)
	row := tr.QueryRowx(getIdQuery, name)
	if err := row.Scan(&itemID, &price); err != nil {
		return 0, err
	}
	getCoinLeft := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userListTable)
	row = tr.QueryRowx(getCoinLeft, userid)
	if err := row.Scan(&amount); err != nil {
		return 0, err
	}
	if amount-price < 0 {
		return 0, errors.New("цена товара выше количества текущих монет")
	}
	createListQuery := fmt.Sprintf("INSERT INTO %s (user_id, item_id, price, purchase_date) VALUES ($1,$2,$3,$4) RETURNING id", purchaseTable)
	row = tr.QueryRowx(createListQuery, userid, itemID, price, time.Now())
	if err := row.Scan(&id); err != nil {
		tr.Rollback()
		return id, err
	}
	changeAmountQuery := fmt.Sprintf("UPDATE %s SET coins = $1 WHERE id = $2", userListTable)
	_, err = tr.Exec(changeAmountQuery, amount-price, userid)
	if err != nil {
		tr.Rollback()
		return 00, err
	}
	logger.Log.Debug().Int("id", id).Msg("Успешно совершена покупка товара")
	return id, tr.Commit()
}

func (r *ShopPostgres) SendCoin(userid int, input domain.Transactions) (int, error) {
	tr, err := r.db.Beginx()

	if err != nil {
		return 0, err
	}
	var id, amount int
	getCoinLeft := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userListTable)
	row := tr.QueryRowx(getCoinLeft, userid)
	if err := row.Scan(&amount); err != nil {
		return 0, err
	}
	if amount-input.Amount < 0 {
		return 0, errors.New("количество отправки выше количества текущих монет")
	}
	sendMoneyQuery := fmt.Sprintf("UPDATE %s SET coins = coins + $1 WHERE id = $2", userListTable)
	row = tr.QueryRowx(sendMoneyQuery, input.Amount, input.Destination)
	if err := row.Scan(&id); err != nil {
		tr.Rollback()
		return id, err
	}
	changeAmountQuery := fmt.Sprintf("UPDATE %s SET coins = coins - $1 WHERE id = $2", userListTable)
	_, err = tr.Exec(changeAmountQuery, input.Amount, userid)
	if err != nil {
		tr.Rollback()
		return 00, err
	}
	createListQuery := fmt.Sprintf("INSERT INTO %s (source, destination, amount, transaction_time) VALUES ($1,$2,$3,$4) RETURNING id", transactionsTable)
	row = tr.QueryRowx(createListQuery, userid, input.Destination, input.Amount, time.Now())
	if err := row.Scan(&id); err != nil {
		tr.Rollback()
		return id, err
	}
	logger.Log.Debug().Int("id", id).Msg("Успешно совершена отправка момент")
	return id, tr.Commit()
}
