package crudlib

import (
	"database/sql"
	"fmt"
)

// DBHandle can be either a *sql.DB or a *sql.Tx. To allow CRUD and hook methods
// to be flexible about transaction management.
type DBHandle interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// InTransaction wraps an operation in BEGIN/COMMIT.
func InTransaction(db *sql.DB, operation func(DBHandle) error) error {

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %s", err)
	}

	err = operation(tx)

	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return fmt.Errorf("rollback failed: %s, rolling back from: %s", err2, err)
		}
		return fmt.Errorf("transaction rolled back: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %s", err)
	}

	return nil
}
