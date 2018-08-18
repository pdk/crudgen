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

## Struct Composition

This system is not smart enough to understand complex structs with composition,
but there is a command-line option to explicitly handle such cases.

For example:

```
-struct User -compose V:Version
```

This indicates that the `User` struct is the one we want to generate CRUD
operations for, and it is composed on a `Version` struct, named `V` within
`User`.

Only a single level of composition is handled, but multiple structs can be
included like `-compose S1:Struct1,S2:Struct2`.

Also, the structs may exist in other files. Just list all the input .go files to
the `crudgen` command.

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

There are Pre- and Post- hooks for insert, update and delete:

1. `PreInsert(tx *sql.TX) error` and `PostInsert(tx *sql.TX) error`
2. `PreUpdate(tx *sql.TX) error` and `PostUpdate(tx *sql.TX) error`
3. `PreDelete(tx *sql.TX) error` and `PostDelete(tx *sql.TX) error`

The current transaction is passed to each hook, so that additional database
operations may be performed.

To use any of these hooks, just define the methods on your struct type. The
methods are detected dynamically via type assertions, so they do not necessarily
need to be in the same source file as your struct. (I.e. the methods are not
detected during code generation, but during execution via type assertions.)

Note that although the Insert() and Update() methods return new struct values,
rather than mutating in place (for ID and timestamps) the hook methods may be
mutating, and the mutations will be returned in the new struct values.

Also, note that the hooks take a transaction handle (`*sql.Tx`) rather than a
database handle (`*sql.DB`). If you have the wrong signature on your hooks, they
will not be found and not invoked.

## Transaction Management

The insert, update and delete operations are each done within a single
transaction, covering their Pre- and Post- hooks. If any error occurs, with the
Pre- or Post- hook, or with the main operation itself, the transaction will be
rolled back.

If the caller wants to include more operations in the transaction, then they
`InsertTx()`, `UpdateTx()` and `DeleteTx()` methods are available. Rather than
passing in a `*sql.DB`, the caller passes in an already active `*sql.Tx`.

The `crudlib.InTransaction()` function is available to wrap a set of operations
in a single transaction.

## Tests

Tests assume a local postgres database to talk to.

To setup test database, do this:

```
create database crud_test;
create user crud_test with password 'MudCrud';
grant all privileges on database crud_test to crud_test;
```
