package sqlxtx

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// SqlxWrapperInterface Wraps sqlx interface
type SqlxWrapperInterface interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, args interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Transactions
}

// Transactions wraps sqlx Transactions functionality
type Transactions interface {
	Rollback() error
	Commit() error
	TxEnd(txFn func() error, config Config) error
}

// Config used as configuration object for EnableSqlxTx
type Config struct {
	Verbose bool
}

// SqlxTxInterface exposes a workable interface for SqlxTx
type SqlxTxInterface interface {
	Exec(txFn func() error) error
}

// EnableSqlxTx init SqlxTx
type EnableSqlxTx struct {
	Client *sqlx.DB
	Config Config
}

func buildSqlxTx(sdb *sqlx.DB, config Config) (SqlxWrapperInterface, error) {
	var sdt SqlxWrapperInterface
	tx, err := sdb.Beginx()
	if err != nil {
		return nil, err
	}
	sdt = &SQLxTx{DB: tx}
	if config.Verbose {
		log.Println("Create Transaction:")
	}
	return sdt, nil
}

// Exec receives a callback meant to swap sqlx.DB by sqlx.Tx in runtime
// it will apply to operation into the callback
func (repo *EnableSqlxTx) Exec(txFn func() error) error {
	Client, _ := buildSqlxTx(repo.Client, repo.Config)
	if err := Client.TxEnd(txFn, repo.Config); err != nil {
		return err
	}
	return nil
}
