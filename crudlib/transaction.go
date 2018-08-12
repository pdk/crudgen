package crudlib

import (
	"database/sql"
	"fmt"
)

// InTransaction wraps an operation in BEGIN/COMMIT.
func InTransaction(db *sql.DB, operation func(*sql.Tx) error) error {

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
