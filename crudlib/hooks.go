package crudlib

// Order of operations with hooks:
// 1. PreInsert, PreUpdate, or PreDelete
// 2. PreModify
// 3. actual SQL insert, update or delete
// 4. PostModify
// 5. PostInsert, PostUpdate, or PostDelete

import "database/sql"

// PreModifier offers an operation to be executed before any
// insert/update/delete operation. Will be executed *after*
// PreInsert/PreUpdate/PreDelete, but *before* actual SQL insert/update/delete.
type PreModifier interface {
	PreModify(*sql.Tx, string) error
}

// PostModifier offers an operation to be executed after any
// insert/update/delete operation. Will be executed *before*
// PostInsert/PostUpdate/PostDelete, but *after* actual SQL
// insert/update/delete.
type PostModifier interface {
	PostModify(*sql.Tx, string) error
}

// PreInserter offers a pre-insert operation which might return an error to
// indicate the operation should be aborted.
type PreInserter interface {
	PreInsert(*sql.Tx, string) error
}

// PostInserter offers a post-insert operation.
type PostInserter interface {
	PostInsert(*sql.Tx, string) error
}

// PreUpdater offers a pre-update operation which might return an error to
// indicate the operation should be aborted.
type PreUpdater interface {
	PreUpdate(*sql.Tx, string) error
}

// PostUpdater offers a post-update operation.
type PostUpdater interface {
	PostUpdate(*sql.Tx, string) error
}

// PreDeleter offers a pre-delete operation.
type PreDeleter interface {
	PreDelete(*sql.Tx, string) error
}

// PostDeleter offers a post-deletion operation.
type PostDeleter interface {
	PostDelete(*sql.Tx, string) error
}

// PreModify checks if the passed in item has a PreModify method and invokes it.
func PreModify(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PreModifier); ok {
		return chk.PreModify(tx, tableName)
	}
	return nil
}

// PostModify checks if the passed in item has a PostModify method and invokes it.
func PostModify(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PostModifier); ok {
		return chk.PostModify(tx, tableName)
	}
	return nil
}

// PreInsert checks if the passed in item has a PreInsert method and invokes it.
func PreInsert(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PreInserter); ok {
		return chk.PreInsert(tx, tableName)
	}
	return nil
}

// PostInsert checks if the passed in item has a PostInsert method and invokes it.
func PostInsert(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PostInserter); ok {
		return chk.PostInsert(tx, tableName)
	}

	return nil
}

// PreUpdate checks if the passed in item has a PreUpdate method and invokes it.
func PreUpdate(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PreUpdater); ok {
		return chk.PreUpdate(tx, tableName)
	}
	return nil
}

// PostUpdate checks if the passed in item has a PostUpdate method and invokes it.
func PostUpdate(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PostUpdater); ok {
		return chk.PostUpdate(tx, tableName)
	}
	return nil
}

// PreDelete checks if the passed in item has a PreDelete method and invokes it.
func PreDelete(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PreDeleter); ok {
		return chk.PreDelete(tx, tableName)
	}
	return nil
}

// PostDelete checks if the passed in item has a PostDelete method and invokes it.
func PostDelete(tx *sql.Tx, item interface{}, tableName string) error {
	if chk, ok := item.(PostDeleter); ok {
		return chk.PostDelete(tx, tableName)
	}
	return nil
}
