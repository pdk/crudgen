// Code generated with github.com/pdk/crudgen DO NOT EDIT.

package samples

import (
    "database/sql"
    "github.com/pdk/crudgen/crudlib"
    "time"
)

// Insert will insert one User instance as a row in table users.
func (r *User) Insert(db *sql.DB) error {
    
    r.V.VersionAt = time.Now()
    
    err := crudlib.PreInsert(r)
    if err != nil {
        return err
    }

    insertStatement := `insert into users (uuid, version_at, active_version, name, email, phone) values ($1, $2, $3, $4, $5, $6) returning version_id`
    
    var newID int64
    err = db.QueryRow(insertStatement, r.V.UUID, r.V.VersionAt, r.V.ActiveVersion, r.Name, r.Email, r.Phone).Scan(&newID)
    r.V.VersionID = newID;
    
	return err
}

// Update will update a row in table users.
func (r *User) Update(db *sql.DB) (rowCount int64, err error) {
    
    r.V.VersionAt = time.Now()
    
    err = crudlib.PreUpdate(r)
    if err != nil {
        return 0, err
    }

    updateStatement := `update users set uuid = $1, version_at = $2, active_version = $3, name = $4, email = $5, phone = $6 where version_id = $7`

    result, err := db.Exec(updateStatement, r.V.UUID, r.V.VersionAt, r.V.ActiveVersion, r.Name, r.Email, r.Phone, r.V.VersionID)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete will delete a row in table users.
func (r *User) Delete(db *sql.DB) (rowCount int64, err error) {

    deleteStatement := `delete from users where version_id = $1`

    result, err := db.Exec(deleteStatement, r.V.VersionID)

    crudlib.PostDelete(r)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
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
