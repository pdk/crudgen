package samples

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

const (
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

	err := x.Insert(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}
}

func TestSelectUsers(t *testing.T) {

	users, err := SelectUsersAll(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	for _, user := range users {
		t.Logf("user: %s\n", user.Email)
	}

}
