package db

import (
	"fmt"
	"log"
	"os"

	"github.com/andresmijares/sqlxtx"
	"github.com/jmoiron/sqlx"
)

var (
	Client *sqlx.DB
	WithTx sqlxtx.SqlxTxInterface // Support for transactions

	username = os.Getenv("username")
	password = os.Getenv("password")
	host     = os.Getenv("host")
	schema   = os.Getenv("schema")

	// GH actions
	githubCI    = "CI"
	githubCiVal = "true"
)

func init() {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host,
		schema)
	var err error

	Client, err = sqlx.Open("mysql", datasourceName)
	if err != nil {
		panic(err)
	}

	WithTx = &sqlxtx.EnableSqlxTx{
		Client: Client,
		Config: sqlxtx.Config{
			Verbose: true,
		},
	}

	if err = Client.Ping(); err != nil {
		if os.Getenv(githubCI) == githubCiVal {
			// if we are testing, ignore this error
			panic(err)
		}
	}

	log.Print("database successfully configured")
}
