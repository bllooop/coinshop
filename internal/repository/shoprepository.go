package repository

import (
	"database/sql"
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
	defer tr.Rollback() // nolint:errcheck

	var id, itemID, price, amount int
	getIdQuery := fmt.Sprintf("SELECT id, price FROM %s WHERE name = $1", shopTable)
	row := tr.QueryRowx(getIdQuery, name)
	if err = row.Scan(&itemID, &price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("товар не найден")
		}

		return 0, err
	}
	getCoinLeft := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userListTable)
	row = tr.QueryRowx(getCoinLeft, userid)
	if err = row.Scan(&amount); err != nil {
		return 0, err
	}
	if amount-price < 0 {
		return 0, errors.New("цена товара выше количества текущих монет")
	}
	createListQuery := fmt.Sprintf("INSERT INTO %s (user_id, item_id, price, purchase_date) VALUES ($1,$2,$3,$4) RETURNING id", purchaseTable)
	row = tr.QueryRowx(createListQuery, userid, itemID, price, time.Now())
	if err = row.Scan(&id); err != nil {
		rollbackErr := tr.Rollback()
		if rollbackErr != nil {
			logger.Log.Error().Err(rollbackErr).Msg("Error during rollback")
		}
		return id, err
	}
	changeAmountQuery := fmt.Sprintf("UPDATE %s SET coins = $1 WHERE id = $2", userListTable)
	_, err = tr.Exec(changeAmountQuery, amount-price, userid)
	if err != nil {
		rollbackErr := tr.Rollback()
		if rollbackErr != nil {
			logger.Log.Error().Err(rollbackErr).Msg("Error during rollback")
		}
		return 0, err
	}
	logger.Log.Debug().Int("id", id).Msg("Успешно совершена покупка товара")
	return id, tr.Commit()
}

func (r *ShopPostgres) SendCoin(input domain.Transactions) (int, error) {
	tr, err := r.beginTransaction()
	if err != nil {
		return 0, err
	}
	defer tr.Rollback() // nolint:errcheck

	amount, err := r.getSourceCoinAmount(tr, *input.Source)
	if err != nil {
		return 0, err
	}

	if amount-input.Amount < 0 {
		return 0, errors.New("количество отправки выше количества текущих монет")
	}

	destId, err := r.getDestinationUserId(tr, input.DestinationUsername)
	if err != nil {
		return 0, err
	}
	input.Destination = &destId
	if err = r.transferCoins(tr, input.Amount, destId, *input.Source); err != nil {
		return 0, err
	}

	id, err := r.createTransaction(tr, input)
	if err != nil {
		return 0, err
	}

	logger.Log.Debug().Int("id", id).Msg("Успешно совершена отправка момент")
	return id, tr.Commit()
}

func (r *ShopPostgres) beginTransaction() (*sqlx.Tx, error) {
	tr, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (r *ShopPostgres) getSourceCoinAmount(tr *sqlx.Tx, userId int) (int, error) {
	var amount int
	getCoinLeft := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userListTable)
	row := tr.QueryRowx(getCoinLeft, userId)
	if err := row.Scan(&amount); err != nil {
		return 0, err
	}
	return amount, nil
}

func (r *ShopPostgres) getDestinationUserId(tr *sqlx.Tx, username string) (int, error) {
	var destId int
	getDestId := fmt.Sprintf("SELECT id FROM %s WHERE username = $1", userListTable)
	row := tr.QueryRowx(getDestId, username)
	if err := row.Scan(&destId); err != nil {
		return 0, err
	}
	return destId, nil
}

func (r *ShopPostgres) transferCoins(tr *sqlx.Tx, amount, destId, sourceId int) error {
	sendMoneyQuery := fmt.Sprintf("UPDATE %s SET coins = coins + $1 WHERE id = $2", userListTable)
	_, err := tr.Exec(sendMoneyQuery, amount, destId)
	if err != nil {
		return err
	}

	changeAmountQuery := fmt.Sprintf("UPDATE %s SET coins = coins - $1 WHERE id = $2", userListTable)
	_, err = tr.Exec(changeAmountQuery, amount, sourceId)
	return err
}

func (r *ShopPostgres) createTransaction(tr *sqlx.Tx, input domain.Transactions) (int, error) {
	createListQuery := fmt.Sprintf("INSERT INTO %s (source, destination, amount, transaction_time) VALUES ($1,$2,$3,$4) RETURNING id", transactionsTable)
	row := tr.QueryRowx(createListQuery, input.Source, input.Destination, input.Amount, input.Timestamp)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
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
    SELECT t.source, u.username AS source_username, t.amount
    FROM transactions t
    JOIN userlist u ON t.source = u.id
    WHERE t.destination = $1;
    `, userID)
	if err != nil {
		return nil, err
	}

	var sentCoins []domain.Transactions
	err = s.db.Select(&sentCoins, `
    SELECT t.destination,d.username AS destination_username, t.amount
    FROM transactions t
    JOIN userlist d ON t.destination = d.id
    WHERE t.source = $1;
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
	logger.Log.Debug().Msg("Успешно совершено получение информации о пользователе")
	return userSummary, nil
}

func (r *ShopPostgres) DB() *sqlx.DB {
	return r.db
}
