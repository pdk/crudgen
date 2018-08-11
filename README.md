# crudgen

Generate CRUD operations for structs.

## Install

```
go get github.com/pdk/crudgen
```

Currently depends on $GOPATH to find the template.

## Example invocation

```
//go:generate crudgen -source $GOFILE -out sample1-crud.go -package $GOPACKAGE -bindstyle dollar -table stories
//go:generate gofmt -w sample1-crud.go
```

(gofmt added cuz default output is not yet canonical.)

## Requirements

The struct fields must be of types suitable for use with database/sql, e.g.
`sql.NullInt64`. There is nothing in the generation process to verify usable
datatypes.

This will work for non-published fields. The generator is not using reflection,
but parses of actual .go source, so non-published fields are included in crud
operations.

On the other hand, composed structs are **not** handled. Only the fields
directly listed in the struct definition are included in crud operations.

## Hooks

There are three available hooks:

1. `PreInsert() error`
2. `PreUpdate() error`
3. `PostDelete()`

The two `Pre`-hooks allow custom operations to be performed before the
insert/update, or validation to be performed. If either returns a non-nil error,
then the database operation will not happen.

The `PostDelete()` hook allows "clean up" operations to be performed after a
record has been deleted. This method has no option to return an error.

To use these hooks, just define the methods on your struct type. The methods are
detected dynamically, so they do not necessarily need to be in the same source
file as your struct. (I.e. the methods are not detected during code generation,
but during execution via type assertions.)

## Tests

Tests assume a local postgres database to talk to.

To setup test database, do this:

```
create database crud_test
create user crud_test with password 'MudCrud';
grant all privileges on database crud_test to crud_test;
```
