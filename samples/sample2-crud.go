// Code generated with github.com/pdk/crudgen DO NOT EDIT.

package samples

import (
	"database/sql"
	"time"

	"github.com/pdk/crudgen/crudlib"
)

// Insert will insert one User instance as a row in table users.
func (r *User) Insert(db *sql.DB) error {
    return crudlib.InTransaction(db, func (tx *sql.Tx) error {
        return r.InsertTx(tx)
    })
}

// InsertTx will insert, given a transaction.
func (r *User) InsertTx(tx *sql.Tx) error {
    
    r.V.VersionAt = time.Now()
    
    err := crudlib.PreInsert(tx, r)
    if err != nil {
        return err
    }

    insertStatement := `insert into users (uuid, version_at, active_version, name, email, phone) values ($1, $2, $3, $4, $5, $6) returning version_id`
    
    var newID int64
    err = tx.QueryRow(insertStatement, r.V.UUID, r.V.VersionAt, r.V.ActiveVersion, r.Name, r.Email, r.Phone).Scan(&newID)
    r.V.VersionID = newID;
    

    if err != nil {
        return err
    }

    return crudlib.PostInsert(tx, r)
}

// Update will update a row in table users.
func (r *User) Update(db *sql.DB) (rowCount int64, err error) {
    err = crudlib.InTransaction(db, func (tx *sql.Tx) error {
        rowCount, err = r.UpdateTx(tx)
        return err
    })

    return rowCount, err
}

// UpdateTx will update a row, within a transaction.
func (r *User) UpdateTx(tx *sql.Tx) (rowCount int64, err error) {
    
    r.V.VersionAt = time.Now()
    
    err = crudlib.PreUpdate(tx, r)
    if err != nil {
        return 0, err
    }

    updateStatement := `update users set uuid = $1, version_at = $2, active_version = $3, name = $4, email = $5, phone = $6 where version_id = $7`

    result, err := tx.Exec(updateStatement, r.V.UUID, r.V.VersionAt, r.V.ActiveVersion, r.Name, r.Email, r.Phone, r.V.VersionID)

	if err != nil {
		return 0, err
	}

    rows, err := result.RowsAffected()
    if err != nil {
        return rows, err
    }

    return rows, crudlib.PostUpdate(tx, r)
}

// Delete will delete a row in table users.
func (r *User) Delete(db *sql.DB) (rowCount int64, err error) {
    err = crudlib.InTransaction(db, func (tx *sql.Tx) error {
        rowCount, err = r.DeleteTx(tx)
        return err
    })

    return rowCount, err
}

// DeleteTx executes PreDelete, delete, and PostDelete within a transaction.
func (r *User) DeleteTx(tx *sql.Tx) (rowCount int64, err error) {

    deleteStatement := `delete from users where version_id = $1`

    err = crudlib.PreDelete(tx, r)
    if err != nil {
        return 0, err
    }

    result, err := tx.Exec(deleteStatement, r.V.VersionID)

	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
    if err != nil {
        return rows, err
    }

    return rows, crudlib.PostDelete(tx, r)
}

// SelectUsers will select records from table users and return a slice of
// User. The additionalClauses argument should be SQL to be
// appended to the "select ... from users" statement, using "?" for bind
// variables.  E.g. "where foo = ?". bindValues must be provided in the correct
// order to match bind placeholders in the additionalClauses.
func SelectUsers(db *sql.DB, additionalClauses string, bindValues ...interface{}) ([]User, error) {

    selectStatement := `select users.uuid, users.version_id, users.version_at, users.active_version, users.name, users.email, users.phone from users`

    if len(additionalClauses) > 0 {
        selectStatement += " " + additionalClauses
        selectStatement = crudlib.DollarNum.Rebind(selectStatement)
    }

    values := []User{}

    rows, err := db.Query(selectStatement, bindValues...)
    if err != nil {
        return values, err
    }
    defer rows.Close()

    for rows.Next() {
        i := User{}
        err := rows.Scan(&i.V.UUID, &i.V.VersionID, &i.V.VersionAt, &i.V.ActiveVersion, &i.Name, &i.Email, &i.Phone)
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

// SelectUsersAll does a Select with no additional conditions/clauses.
func SelectUsersAll(db *sql.DB) ([]User, error) {
    return SelectUsers(db, "")
}

// SelectUsersRow will select one record from table users and return a
// User. The additionalClauses argument should be SQL to be
// appended to the "select ... from users" statement, using "?" for bind
// variables.  E.g. "where foo = ?". bindValues must be provided in the correct
// order to match bind placeholders in the additionalClauses.
// Returns sql.ErrNoRows if no rows found.
func SelectUsersRow(db *sql.DB, additionalClauses string, bindValues ...interface{}) (User, error) {

    selectStatement := `select users.uuid, users.version_id, users.version_at, users.active_version, users.name, users.email, users.phone from users`

    if len(additionalClauses) > 0 {
        selectStatement += " " + additionalClauses
        selectStatement = crudlib.DollarNum.Rebind(selectStatement)
    }

    i := User{}

    err := db.QueryRow(selectStatement, bindValues...).Scan(
        &i.V.UUID, &i.V.VersionID, &i.V.VersionAt, &i.V.ActiveVersion, &i.Name, &i.Email, &i.Phone)

    return i, err
}
