package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
