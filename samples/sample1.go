package samples

//go:generate crudgen -source $GOFILE -out sample1-crud.go -package $GOPACKAGE -bindstyle dollar -table stories
//go:generate gofmt -w sample1-crud.go

import (
	"database/sql"
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
