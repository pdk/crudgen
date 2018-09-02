package user

//go:generate crudgen -out user-crud.go -package $GOPACKAGE -bindstyle dollar -table users -select SelectUsers -struct User -compose V:Version $GOFILE ../version/version.go
//go:generate gofmt -w user-crud.go

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pdk/crudgen/samples/version"
)

// User is a human who visits our site.
type User struct {
	V     version.Version
	Name  string
	Email string
	Phone string
}

func (u *User) PreInsert(tx *sql.Tx, tableName string) error {
	if u.V.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		u.V.UUID = uuid.New()
	}

	return nil
}

func (u *User) PostInsert(tx *sql.Tx, tableName string) error {

	deactivateSQL := `
		update users
		set active_version = false
		where active_version = true
		and uuid = $1`

	activateSQL := `
		update users
		set active_version = true
		where active_version = false
		and uuid = $1
		and version_id = (
			select max(version_id)
			from users
			where uuid = $2
		)`

	_, err := tx.Exec(deactivateSQL, u.V.UUID)
	if err != nil {
		return fmt.Errorf("unable to deactivate old versions for %s: %s", u.V.UUID, err)
	}

	_, err = tx.Exec(activateSQL, u.V.UUID, u.V.UUID)
	if err != nil {
		return fmt.Errorf("unable to activate newest versions for %s: %s", u.V.UUID, err)
	}

	return nil
}
