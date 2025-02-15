package repository

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestShopPostgres_BuyItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	type args struct {
		userid int
		name   string
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id, price FROM shop WHERE name = (.+)").
					WithArgs("cup").
					WillReturnRows(sqlmock.NewRows([]string{"id", "price"}).AddRow(1, 10))

				mock.ExpectQuery("SELECT coins FROM userlist WHERE id = (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(100))

				mock.ExpectQuery("INSERT INTO purchases").
					WithArgs(1, 1, 10, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectExec("UPDATE userlist SET coins = (.+) WHERE id = (.+)").
					WithArgs(90, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			input: args{
				userid: 1,
				name:   "cup",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error during execution (Insufficient funds)",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id, price FROM shop WHERE name = (.+)").
					WithArgs("cup").
					WillReturnRows(sqlmock.NewRows([]string{"id", "price"}).AddRow(1, 10))

				mock.ExpectQuery("SELECT coins FROM userlist WHERE id = (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(5))

				mock.ExpectRollback()
			},
			input: args{
				userid: 1,
				name:   "cup",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Error during execution (Item not found)",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id, price FROM shop WHERE name = (.+)").
					WithArgs("cup").
					WillReturnError(errors.New("item not found"))

				mock.ExpectRollback()
			},
			input: args{
				userid: 1,
				name:   "cup",
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			shop := NewShopPostgres(sqlxDB)

			got, err := shop.BuyItem(tt.input.userid, tt.input.name)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestShopPostgres_sendCoin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	tests := []struct {
		name    string
		mock    func()
		input   domain.Transactions
		want    int
		wantErr bool
	}{

		{
			name: "OK",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(fmt.Sprintf("SELECT coins FROM %s (.+)", userListTable)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s WHERE (.+)", userListTable)).
					WithArgs("name").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+)", userListTable)).
					WithArgs(10, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+)", userListTable)).
					WithArgs(10, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s (.+)", transactionsTable)).
					WithArgs(1, 2, 10, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			input: domain.Transactions{
				Source:              IntPointer(1),
				DestinationUsername: "name",
				Destination:         IntPointer(2),
				Amount:              10,
				Timestamp:           func() *time.Time { t := time.Now(); return &t }(),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error during execution (Insufficient funds)",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(fmt.Sprintf("SELECT coins FROM %s WHERE id = (.+)", userListTable)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(10))
			},
			input: domain.Transactions{
				Source:              IntPointer(1),
				DestinationUsername: "name",
				Destination:         IntPointer(2),
				Amount:              10,
				Timestamp:           func() *time.Time { t := time.Now(); return &t }(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Error - Insert Transaction Fails",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(fmt.Sprintf("SELECT coins FROM %s (.+)", userListTable)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(100))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s WHERE (.+)", userListTable)).
					WithArgs("name").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+)", userListTable)).
					WithArgs(10, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+)", userListTable)).
					WithArgs(10, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(fmt.Sprintf(`INSERT INTO %s (.+)`, transactionsTable)).
					WithArgs(1, 2, 10, sqlmock.AnyArg()).
					WillReturnError(errors.New("insert failed"))

				mock.ExpectRollback()
			},
			input: domain.Transactions{
				Source:              IntPointer(1),
				Destination:         IntPointer(2),
				DestinationUsername: "name",
				Amount:              10,
				Timestamp:           func() *time.Time { t := time.Now(); return &t }(),
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			shop := NewShopPostgres(sqlxDB)

			got, err := shop.SendCoin(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
