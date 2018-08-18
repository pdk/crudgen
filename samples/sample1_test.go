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
	host                        = "localhost"
	port                        = 5432
	user                        = "crud_test"
	password                    = "MudCrud"
	dbName                      = "crud_test"
	createStoriesTableStatement = `
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

	_, err = globalDB.Exec(createStoriesTableStatement)
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
	}

	baseTime := time.Now()

	x, err := x.Insert(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if x.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to have a value")
	}

	if x.CreatedAt.Before(baseTime) {
		t.Errorf("expected CreatedAt to have a time later than %s, but got %s", baseTime.String(), x.CreatedAt.String())
	}

	data, err := SelectAll(globalDB)

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

	// verify the PreInsert method was invoked.
	if data[0].URL != "https://www.guam.net" {
		t.Errorf("expected URL to be https://www.guam.net, but got %s", data[0].URL)
	}

	x.URL = "http://google.com/"
	x.Name = ""

	timeTwo := time.Now()

	x, err = x.Update(globalDB)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	data1, err := SelectRow(globalDB, "where id = ?", x.ID)

	if err != nil {
		t.Errorf("did not expect error, but got %s", err)
	}

	if data1.URL != x.URL {
		t.Errorf("expected URL to be %s, but got %s", x.URL, data1.URL)
	}

	if x.UpdatedAt.Before(timeTwo) {
		t.Errorf("expected updated_at to be updated, but %s is before %s", x.UpdatedAt.String(), timeTwo.String())
	}

	if data1.Name != "(no name available)" {
		t.Errorf("expected name to be '(no name available)', but got '%s'", data1.Name)
	}

	rowCount, err := x.Delete(globalDB)

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
