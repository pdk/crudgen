// Code generated with github.com/pdk/crudgen DO NOT EDIT.

package story

import (
	"database/sql"
	"time"

	"github.com/pdk/crudgen/crudlib"
)

const (
	TableName = "stories"
)

// Insert wraps InsertTx() in a transaction.
func (r Story) Insert(db *sql.DB) (newR Story, err error) {

	err = crudlib.InTransaction(db, func(tx *sql.Tx) error {
		newR, err = r.InsertTx(tx)
		return err
	})

	return newR, err
}

// InsertTx will insert on Story, given a transaction. Invokes
// PreInsert and PostInsert hooks. Returns new Story with ID,
// timestamps, etc updated.
func (r Story) InsertTx(tx *sql.Tx) (Story, error) {

	err := crudlib.PreInsert(tx, &r, TableName)
	if err != nil {
		return r, err
	}
	err = crudlib.PreModify(tx, &r, TableName)
	if err != nil {
		return r, err
	}

	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()

	insertStatement := `insert into stories (url, mp3_url, mp3_duration, image_urls, name, description, place, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var newID int64
	err = tx.QueryRow(insertStatement, r.URL, r.MP3URL, r.MP3Duration, r.imageURLs, r.Name, r.Description, r.place, r.CreatedAt, r.UpdatedAt).Scan(&newID)
	r.ID = newID

	if err != nil {
		return r, err
	}

	err = crudlib.PostModify(tx, &r, TableName)
	if err != nil {
		return r, err
	}
	err = crudlib.PostInsert(tx, &r, TableName)

	return r, err
}

// Wraps UpdateTx() in a transaction.
func (r Story) Update(db *sql.DB) (newR Story, err error) {

	err = crudlib.InTransaction(db, func(tx *sql.Tx) error {
		newR, err = r.UpdateTx(tx)
		return err
	})

	return newR, err
}

// UpdateTx will update a row, given a transaction. Will fail and return an
// error if not exactly 1 row is updated. Returns new Story with
// timestamps, etc updated.
func (r Story) UpdateTx(tx *sql.Tx) (Story, error) {

	err := crudlib.PreUpdate(tx, &r, TableName)
	if err != nil {
		return r, err
	}
	err = crudlib.PreModify(tx, &r, TableName)
	if err != nil {
		return r, err
	}

	r.UpdatedAt = time.Now()

	updateStatement := `update stories set url = $1, mp3_url = $2, mp3_duration = $3, image_urls = $4, name = $5, description = $6, place = $7, created_at = $8, updated_at = $9 where id = $10`

	result, err := tx.Exec(updateStatement, r.URL, r.MP3URL, r.MP3Duration, r.imageURLs, r.Name, r.Description, r.place, r.CreatedAt, r.UpdatedAt, r.ID)

	rows, err := result.RowsAffected()
	if err != nil {
		return r, err
	}

	if rows == 0 {
		return r, crudlib.NoRowsUpdated
	}

	if rows > 1 {
		return r, crudlib.MoreThanOneRowUpdated
	}

	err = crudlib.PostModify(tx, &r, TableName)
	if err != nil {
		return r, err
	}
	err = crudlib.PostUpdate(tx, &r, TableName)

	return r, err
}

// Wraps DeleteTx() in a transaction.
func (r *Story) Delete(db *sql.DB) (rowCount int64, err error) {
	err = crudlib.InTransaction(db, func(tx *sql.Tx) error {
		rowCount, err = r.DeleteTx(tx)
		return err
	})

	return rowCount, err
}

// DeleteTx executes PreDelete, delete, and PostDelete within a transaction.
// Will fail if delete affects more than one row. No error if 0 rows are
// deleted.
func (r *Story) DeleteTx(tx *sql.Tx) (rowCount int64, err error) {

	deleteStatement := `delete from stories where id = $1`

	err = crudlib.PreDelete(tx, r, TableName)
	if err != nil {
		return 0, err
	}
	err = crudlib.PreModify(tx, &r, TableName)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(deleteStatement, r.ID)

	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return rows, err
	}

	if rows > 1 {
		return rows, crudlib.MoreThanOneRowDeleted
	}

	err = crudlib.PostModify(tx, &r, TableName)
	if err != nil {
		return rows, err
	}
	return rows, crudlib.PostDelete(tx, r, TableName)
}

// Select wraps SelectTx in a transaction.
func Select(db *sql.DB, additionalClauses string, bindValues ...interface{}) (results []Story, err error) {

	err = crudlib.InTransaction(db, func(tx *sql.Tx) error {
		results, err = SelectTx(tx, additionalClauses, bindValues...)
		return err
	})

	return results, err
}

// SelectTx will select records from table stories and return a slice of
// Story. The additionalClauses argument should be SQL to be
// appended to the "select ... from stories" statement, using "?" for bind
// variables.  E.g. "where foo = ?". bindValues must be provided in the correct
// order to match bind placeholders in the additionalClauses.
func SelectTx(tx *sql.Tx, additionalClauses string, bindValues ...interface{}) ([]Story, error) {

	selectStatement := `select stories.id, stories.url, stories.mp3_url, stories.mp3_duration, stories.image_urls, stories.name, stories.description, stories.place, stories.created_at, stories.updated_at from stories`

	if len(additionalClauses) > 0 {
		selectStatement += " " + additionalClauses
		selectStatement = crudlib.DollarNum.Rebind(selectStatement)
	}

	values := []Story{}

	rows, err := tx.Query(selectStatement, bindValues...)
	if err != nil {
		return values, err
	}
	defer rows.Close()

	for rows.Next() {
		i := Story{}
		err := rows.Scan(&i.ID, &i.URL, &i.MP3URL, &i.MP3Duration, &i.imageURLs, &i.Name, &i.Description, &i.place, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return values, err
		}
		values = append(values, i)
	}

	err = rows.Err()
	if err != nil {
		return values, err
	}

	return values, rows.Close()
}

// SelectAll does a Select with no additional conditions/clauses.
func SelectAll(db *sql.DB) ([]Story, error) {
	return Select(db, "")
}

// SelectRow will select one record from table stories and return a
// Story. The additionalClauses argument should be SQL to be
// appended to the "select ... from stories" statement, using "?" for bind
// variables.  E.g. "where foo = ?". bindValues must be provided in the correct
// order to match bind placeholders in the additionalClauses.
// Returns sql.ErrNoRows if no rows found.
func SelectRow(db *sql.DB, additionalClauses string, bindValues ...interface{}) (Story, error) {

	selectStatement := `select stories.id, stories.url, stories.mp3_url, stories.mp3_duration, stories.image_urls, stories.name, stories.description, stories.place, stories.created_at, stories.updated_at from stories`

	if len(additionalClauses) > 0 {
		selectStatement += " " + additionalClauses
		selectStatement = crudlib.DollarNum.Rebind(selectStatement)
	}

	i := Story{}

	err := db.QueryRow(selectStatement, bindValues...).Scan(
		&i.ID, &i.URL, &i.MP3URL, &i.MP3Duration, &i.imageURLs, &i.Name, &i.Description, &i.place, &i.CreatedAt, &i.UpdatedAt)

	return i, err
}
