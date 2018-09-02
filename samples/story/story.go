package story

//go:generate crudgen -out story-crud.go -package $GOPACKAGE -bindstyle dollar -table stories -select Select $GOFILE
//go:generate gofmt -w story-crud.go

import (
	"database/sql"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Story has details about a single story.
type Story struct {
	ID          int64          `db:"id" crud:"autoincr"`
	URL         string         `db:"url"`
	MP3URL      string         `db:"mp3_url"`
	MP3Duration sql.NullInt64  `db:"mp3_duration"`
	imageURLs   pq.StringArray `db:"image_urls"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	place       string         `db:"place"`
	CreatedAt   time.Time      `db:"created_at" crud:"create_timestamp"`
	UpdatedAt   time.Time      `db:"updated_at" crud:"update_timestamp"`
}

// PreInsert runs before inserting a row.
func (s *Story) PreInsert(tx *sql.Tx, tableName string) error {
	if strings.HasPrefix(s.URL, "http://") {
		s.URL = strings.Replace(s.URL, "http://", "https://", 1)
	}

	return nil
}

// PreUpdate runs before updating a row.
func (s *Story) PreUpdate(tx *sql.Tx, tableName string) error {
	if s.Name == "" {
		s.Name = "(no name available)"
	}

	return nil
}
