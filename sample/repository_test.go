package sample

import (
	"errors"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func NewMock() (*sqlx.DB, sqlxmock.Sqlmock) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCreateRepoError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectExec(`INSERT INTO users \(first_name, last_name\) VALUES \(\?, \?\);`).
		WillReturnError(errors.New("database error"))

	err := repo.Create()

	assert.NotNil(t, err)
	assert.EqualValues(t, err.Error(), "database error")
}

func TestCreateRepoNoError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	repo := NewUserRepository(db)

	prep := mock.ExpectExec(`INSERT INTO users \(first_name, last_name\) VALUES \(\?, \?\);`)
	prep.WithArgs(
		"Jhon",
		"Doe",
	).WillReturnResult(sqlxmock.NewResult(1, 1))

	err := repo.Create()

	assert.Nil(t, err)
}

func TestEventRepoError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	repo := NewEventRepository(db)

	mock.ExpectExec(`INSERT INTO events\(name\) VALUE \(\?\);`).
		WillReturnError(errors.New("database error"))

	err := repo.Create("create_user")

	assert.NotNil(t, err)
	assert.EqualValues(t, err.Error(), "database error")
}

func TestEventRepoNoError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	repo := NewEventRepository(db)

	prep := mock.ExpectExec(`INSERT INTO events\(name\) VALUE \(\?\);`)
	prep.WithArgs(
		"create_user",
	).WillReturnResult(sqlxmock.NewResult(1, 1))

	err := repo.Create("create_user")

	assert.Nil(t, err)
}
