package sample

import (
	"errors"
	"testing"

	"github.com/andresmijares/sqlxtx/sample/db"
	"github.com/stretchr/testify/assert"
)

var (
	getUserMock  func() error
	getEventMock func() error
	getExecMock  func(func() error) error
)

// ******Mocks*******
type userDaoMock struct{}

func (u *userDaoMock) Create() error {
	return getUserMock()
}

type eventDaoMock struct{}

func (u *eventDaoMock) Create(name string) error {
	return getEventMock()
}

type WithTxMock struct{}

func (repo *WithTxMock) Exec(fn func() error) error {
	return getExecMock(fn)
}

func init() {
	db.WithTx = &WithTxMock{}
	UserDao = &userDaoMock{}
	EventsDao = &eventDaoMock{}
}

func TestCreateUserDaoCreateError(t *testing.T) {
	getUserMock = func() error {
		return errors.New("database error (user)")
	}
	getEventMock = func() error {
		return nil
	}

	err := UsersServices.Create()
	assert.NotNil(t, err)
	assert.EqualValues(t, err.Error(), "database error (user)")
}

func TestCreateEventDaoCreateError(t *testing.T) {
	getUserMock = func() error {
		return nil
	}
	getEventMock = func() error {
		return errors.New("database error (event)")
	}

	err := UsersServices.Create()
	assert.NotNil(t, err)
	assert.EqualValues(t, err.Error(), "database error (event)")
}

func TestCreateNoError(t *testing.T) {
	getUserMock = func() error {
		return nil
	}
	getEventMock = func() error {
		return nil
	}

	err := UsersServices.Create()
	assert.Nil(t, err)
}

func TestCreateWithTxEventUserDaoCreateError(t *testing.T) {
	var TxErr error
	getExecMock = func(fn func() error) error {
		TxErr = fn()
		return errors.New("database error (user)")
	}
	getUserMock = func() error {
		return errors.New("database error (user)")
	}
	getEventMock = func() error {
		return nil
	}

	err := UsersServices.CreateWithTx()
	assert.NotNil(t, err)
	assert.NotNil(t, TxErr)
	assert.EqualValues(t, err.Error(), TxErr.Error())
}

func TestCreateWithTxEventDaoCreateError(t *testing.T) {
	var TxErr error
	getExecMock = func(fn func() error) error {
		TxErr = fn()
		return errors.New("database error (events)")
	}
	getUserMock = func() error {
		return nil
	}
	getEventMock = func() error {
		return errors.New("database error (events)")
	}

	err := UsersServices.CreateWithTx()
	assert.NotNil(t, err)
	assert.NotNil(t, TxErr)
	assert.EqualValues(t, err.Error(), TxErr.Error())
}

func TestCreateWithTxTransactionError(t *testing.T) {
	var TxErr error
	getExecMock = func(fn func() error) error {
		TxErr = fn()
		return nil
	}
	getUserMock = func() error {
		return nil
	}
	getEventMock = func() error {
		return nil
	}

	err := UsersServices.CreateWithTx()
	assert.Nil(t, err)
	assert.Nil(t, TxErr)
}

func TestCreateWithTxCreateNoError(t *testing.T) {

}
