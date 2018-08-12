package samples

//go:generate crudgen -out sample2-crud.go -package $GOPACKAGE -bindstyle dollar -table users -select SelectUsers -struct User -compose V:Version $GOFILE ./version/version.go
//go:generate gofmt -w sample1-crud.go

import (
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

func (u *User) PreInsert() error {
	u.V.UUID = uuid.New()

	return nil
}
