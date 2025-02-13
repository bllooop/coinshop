package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bllooop/coinshop/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	r := NewAuthPostgres(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   domain.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO userlist").
					WithArgs("username", "123", 1000).WillReturnRows(rows)
			},
			input: domain.User{
				UserName: "username",
				Password: "123",
				Coins:    intPointer(1000),
			},
			want: 1,
		},
		{
			name: "Empty input fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO userlist").
					WithArgs("", "123", 1000).WillReturnRows(rows)
			},
			input: domain.User{
				UserName: "",
				Password: "123",
				Coins:    intPointer(1000),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
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

func TestAuthPostgres_SignUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	r := NewAuthPostgres(sqlxDB)

	type args struct {
		username string
		password string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    domain.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, "test", "password")
				mock.ExpectQuery("SELECT (.+) FROM userlist").
					WithArgs("test", "password").WillReturnRows(rows)
			},
			input: args{"test", "password"},
			want: domain.User{
				Id:       1,
				UserName: "test",
				Password: "password",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "coins"})
				mock.ExpectQuery("SELECT (.+) FROM userlist").
					WithArgs("not", "found").WillReturnRows(rows)
			},
			input:   args{"not", "found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.SignUser(tt.input.username, tt.input.password)
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

func intPointer(s int) *int {
	return &s
}
