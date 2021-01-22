package sample

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepoInterface interface {
	Create() error
}

type userRepository struct {
	Client *sqlx.DB
}

func (u *userRepository) Create() error {
	_, err := u.Client.NamedExec("INSERT INTO users (first_name, last_name) VALUES (:first_name, :last_name);",
		map[string]interface{}{
			"first_name": "Jhon",
			"last_name":  "Doe",
		})
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *sqlx.DB) UserRepoInterface {
	return &userRepository{
		Client: db,
	}
}

type EventRepoInterface interface {
	Create(string) error
}

type eventsRepository struct {
	Client *sqlx.DB
}

func (e *eventsRepository) Create(name string) error {
	_, err := e.Client.NamedExec(`INSERT INTO events(name) VALUE (:name);`,
		map[string]interface{}{
			"name": name,
		})
	if err != nil {
		return err
	}
	return nil
}

func NewEventRepository(db *sqlx.DB) EventRepoInterface {
	return &eventsRepository{
		Client: db,
	}
}
