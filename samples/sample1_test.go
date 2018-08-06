package samples

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	host                 = "localhost"
	port                 = 5432
	user                 = "crud_test"
	password             = "MudCrud"
	dbName               = "crud_test"
	createTableStatement = `
		create table if not exists stories (
			id serial not null primary key,
			url text,
			mp3_url text,
			mp3_duration integer,
			image_urls text[],
			name text not null,
			description text not null,
			place text,
			created_at timestamp without time zone not null,
			updated_at timestamp without time zone not null
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

	_, err = globalDB.Exec(createTableStatement)
	if err != nil {
		log.Fatalf("%s", err)
	}

}

func TestCrud(t *testing.T) {

	x := Story{
		Name:        "a name",
		Description: "the descriptionn",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := x.Insert(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	x.URL = "http://google.com/"

	rowCount, err := x.Update(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if rowCount != 1 {
		t.Errorf("expected to update 1 row, but got %d", rowCount)
	}

	rowCount, err = x.Delete(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if rowCount != 1 {
		t.Errorf("expected to delete 1 row, but got %d", rowCount)
	}
}
