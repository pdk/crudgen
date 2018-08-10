package samples

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/lib/pq"
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
		URL:         "http://www.guam.net",
		MP3URL:      "http://s3.aws.com/soundfile.mp3",
		MP3Duration: sql.NullInt64{Valid: true, Int64: 8762363},
		imageURLs:   pq.StringArray([]string{"http://imgur.com/aoefijowaefj"}),
		Name:        "a name",
		Description: "the description",
		place:       "over there",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := x.Insert(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	data, err := Select(globalDB, "")

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if len(data) == 0 {
		t.Errorf("expected to find rows, but got 0")
	}

	data, err = Select(globalDB, "where id = ?", x.ID)

	if len(data) != 1 {
		t.Errorf("expected to find 1 row, but got %d", len(data))
	}

	if data[0].URL != x.URL {
		t.Errorf("expected URL to be %s, but got %s", x.URL, data[0].URL)
	}

	x.URL = "http://google.com/"

	rowCount, err := x.Update(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if rowCount != 1 {
		t.Errorf("expected to update 1 row, but got %d", rowCount)
	}

	data, err = Select(globalDB, "where id = ?", x.ID)

	if len(data) != 1 {
		t.Errorf("expected to find 1 row, but got %d", len(data))
	}

	if data[0].URL != x.URL {
		t.Errorf("expected URL to be %s, but got %s", x.URL, data[0].URL)
	}

	rowCount, err = x.Delete(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if rowCount != 1 {
		t.Errorf("expected to delete 1 row, but got %d", rowCount)
	}

	data, err = Select(globalDB, "where id = ?", x.ID)

	if len(data) != 0 {
		t.Errorf("expected to find no row, but got %d", len(data))
	}
}
