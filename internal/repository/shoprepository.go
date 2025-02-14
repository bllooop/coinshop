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
		return 0, err
	}
	logger.Log.Debug().Int("id", id).Msg("Успешно совершена покупка товара")
	return id, tr.Commit()
}

func (r *ShopPostgres) SendCoin(input domain.Transactions) (int, error) {
	tr, err := r.db.Beginx()

	if err != nil {
		return 0, err
	}
	var id, amount int
	getCoinLeft := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userListTable)
	row := tr.QueryRowx(getCoinLeft, input.Source)
	if err := row.Scan(&amount); err != nil {
		return 0, err
	}
	if amount-input.Amount < 0 {
		return 0, errors.New("количество отправки выше количества текущих монет")
	}
	sendMoneyQuery := fmt.Sprintf("UPDATE %s SET coins = coins + $1 WHERE username = $2", userListTable)
	_, err = tr.Exec(sendMoneyQuery, input.Amount, input.Destination)
	if err != nil {
		tr.Rollback()
		return 0, err
	}
	changeAmountQuery := fmt.Sprintf("UPDATE %s SET coins = coins - $1 WHERE id = $2", userListTable)
	_, err = tr.Exec(changeAmountQuery, input.Amount, input.Source)
	if err != nil {
		tr.Rollback()
		return 0, err
	}
	createListQuery := fmt.Sprintf("INSERT INTO %s (source, destination, amount, transaction_time) VALUES ($1,$2,$3,$4) RETURNING id", transactionsTable)
	row = tr.QueryRowx(createListQuery, input.Source, input.Destination, input.Amount, input.Timestamp)
	if err := row.Scan(&id); err != nil {
		tr.Rollback()
		return id, err
	}
	logger.Log.Debug().Int("id", id).Msg("Успешно совершена отправка момент")
	return id, tr.Commit()
}

func (s *ShopPostgres) GetUserSummary(userID int) (*domain.UserSummary, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT id, username, coins FROM userlist WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	var purchases []domain.PurchasedItem
	err = s.db.Select(&purchases, `
    SELECT s.name AS item_name, COUNT(*) AS quantity
              FROM purchases p
              JOIN shop s ON p.item_id = s.id
              WHERE p.user_id = $1
              GROUP BY s.name;
    `, userID)
	if err != nil {
		return nil, err
	}

	var receivedCoins []domain.Transactions
	err = s.db.Select(&receivedCoins, `
        SELECT t.source, u.username AS destination, t.amount
        FROM transactions t
        JOIN userlist u ON t.source = u.id
        WHERE t.destination = $1
    `, userID)
	if err != nil {
		return nil, err
	}

	var sentCoins []domain.Transactions
	err = s.db.Select(&sentCoins, `
        SELECT t.destination, u.username AS source, t.amount
        FROM transactions t
        JOIN userlist u ON t.destination = u.id
        WHERE t.source = $1
    `, userID)
	if err != nil {
		return nil, err
	}

	userSummary := &domain.UserSummary{
		UserName:       user.UserName,
		Coins:          *user.Coins,
		PurchasedItems: purchases,
		TransactionsSummary: domain.TransactionsSummary{
			ReceivedCoins: receivedCoins,
			SentCoins:     sentCoins,
		},
	}

	return userSummary, nil
}
