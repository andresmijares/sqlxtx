package sample

import "github.com/andresmijares/sqlxtx/sample/db"

var (
	// UsersServices ~
	UsersServices usersServicesInterface
	UserDao       UserRepoInterface
	EventsDao     EventRepoInterface
)

type usersServices struct{}

type usersServicesInterface interface {
	CreateWithTx() error
	Create() error
}

func init() {
	UserDao = NewUserRepository(db.Client)
	EventsDao = NewEventRepository(db.Client)
	UsersServices = &usersServices{}
}

// Independent operations, either can fail without affecting the other one
func (s *usersServices) Create() error {
	if err := UserDao.Create(); err != nil {
		return err
	}

	if err := EventsDao.Create("userCreated"); err != nil {
		return err
	}

	return nil
}

// ACID operations, if one fail, all fail
func (s *usersServices) CreateWithTx() error {
	if err := db.WithTx.Exec(func() error {
		if err := UserDao.Create(); err != nil {
			return err
		}

		if err := EventsDao.Create("userCreated"); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
