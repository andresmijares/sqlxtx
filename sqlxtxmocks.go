package sqlxtx

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// SQLxTx is wrapper around sqlx.Tx, to change behavior in runtime between sqlx.DB and sqlx.Tx
type SQLxTx struct {
	DB *sqlx.Tx
}

// Select within a transaction.
func (sdt *SQLxTx) Select(dest interface{}, query string, args ...interface{}) error {
	return sdt.DB.Select(dest, query, args...)
}

// Get within a transaction.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (sdt *SQLxTx) Get(dest interface{}, query string, args ...interface{}) error {
	return sdt.DB.Get(dest, query, args...)
}

// NamedExec within a transaction.
// Any named placeholder parameters are replaced with fields from arg.
func (sdt *SQLxTx) NamedExec(query string, args interface{}) (sql.Result, error) {
	return sdt.DB.NamedExec(query, args)
}

// Exec executes a query that doesn't return rows.
// For example: an INSERT and UPDATE.
func (sdt *SQLxTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return sdt.DB.Exec(query, args...)
}

// Rollback aborts the transaction.
func (sdt *SQLxTx) Rollback() error {
	return sdt.DB.Rollback()
}

// Commit commits the transaction.
func (sdt *SQLxTx) Commit() error {
	return sdt.DB.Commit()
}

// TxEnd Finish a transaction
func (sdt *SQLxTx) TxEnd(txFn func() error, config Config) error {
	var err error
	tx := sdt
	defer func() {
		if p := recover(); p != nil {
			if config.Verbose {
				log.Println("Transaction Error, Exec Rollback:", p)
			}
			tx.Rollback()
		} else if err != nil {
			if config.Verbose {
				log.Println("Transaction Error, Exec Rollback:", err)
			}
			tx.Rollback()
		} else {
			if config.Verbose {
				log.Println("Transaction Commited:")
			}
			err = tx.Commit()
		}
	}()

	err = txFn()
	return err
}
