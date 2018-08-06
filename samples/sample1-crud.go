// Code generated with github.com/pdk/crudgen DO NOT EDIT.

package samples

import (
	"database/sql"
)

// Insert will insert one Story instance as a row in table stories.
func (r *Story) Insert(db *sql.DB) error {

	insertStatement := `insert into stories (url, mp3_url, mp3_duration, image_urls, name, description, place, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var newID int64
	err := db.QueryRow(insertStatement, r.URL, r.MP3URL, r.MP3Duration, r.imageURLs, r.Name, r.Description, r.place, r.CreatedAt, r.UpdatedAt).Scan(&newID)
	r.ID = newID

	return err
}

// Update will update a row in table stories.
func (r *Story) Update(db *sql.DB) (rowCount int64, err error) {

	updateStatement := `update stories set url = $1, mp3_url = $2, mp3_duration = $3, image_urls = $4, name = $5, description = $6, place = $7, created_at = $8, updated_at = $9 where id = $10`

	result, err := db.Exec(updateStatement, r.URL, r.MP3URL, r.MP3Duration, r.imageURLs, r.Name, r.Description, r.place, r.CreatedAt, r.UpdatedAt, r.ID)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete will delete a row in table stories.
func (r *Story) Delete(db *sql.DB) (rowCount int64, err error) {

	deleteStatement := `delete from stories where id = $1`

	result, err := db.Exec(deleteStatement, r.ID)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
