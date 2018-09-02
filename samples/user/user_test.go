package user

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "crud_test"
	password = "MudCrud"
	dbName   = "crud_test"

	createUsersTableStatement = `
		create table if not exists users (
			uuid			varchar,
			version_id		serial not null primary key,
			version_at		timestamp,
			active_version 	boolean not null default false,
			name			varchar,
			email			varchar,
			phone			varchar
		)`
)

var globalDB *sql.DB

func init() {
	dbConnectString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	var err error
	globalDB, err = sql.Open("postgres", dbConnectString)
	if err != nil {
		log.Fatalf("%s", err)
	}

	_, err = globalDB.Exec(createUsersTableStatement)
	if err != nil {
		log.Fatalf("%s", err)
	}

}

func TestUsersCrud(t *testing.T) {
	x := User{
		Name:  "George Jetson",
		Email: "george@flybynite.com",
		Phone: "666-999-1234",
	}

	oldUUID := x.V.UUID

	x, err := x.Insert(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if x.V.UUID == oldUUID {
		t.Errorf("expected to get a new UUID")
	}

	oldUUID = x.V.UUID
	x.Phone = "999-666-1234"

	x, err = x.Insert(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if x.V.UUID != oldUUID {
		t.Errorf("did not expect UUID to change")
	}
}

func TestSelectUsers(t *testing.T) {

	users, err := SelectUsersAll(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	for _, user := range users {
		t.Logf("%s %s\n", user.V.UUID, user.Email)
	}

	_, _ = globalDB.Exec("delete from users")
}
