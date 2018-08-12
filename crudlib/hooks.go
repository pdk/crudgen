package crudlib

// PreInserter offers a pre-insert operation which might return an error to
// indicate the operation should be aborted.
type PreInserter interface {
	PreInsert(DBHandle) error
}

// PostInserter offers a post-insert operation.
type PostInserter interface {
	PostInsert(DBHandle) error
}

// PreUpdater offers a pre-update operation which might return an error to
// indicate the operation should be aborted.
type PreUpdater interface {
	PreUpdate(DBHandle) error
}

// PostUpdater offers a post-update operation.
type PostUpdater interface {
	PostUpdate(DBHandle) error
}

// PreDeleter offers a pre-delete operation.
type PreDeleter interface {
	PreDelete(DBHandle) error
}

// PostDeleter offers a post-deletion operation.
type PostDeleter interface {
	PostDelete(DBHandle) error
}

// PreInsert checks if the passed in item has a PreInsert method and invokes it.
func PreInsert(db DBHandle, item interface{}) error {
	if chk, ok := item.(PreInserter); ok {
		return chk.PreInsert(db)
	}
	return nil
}

// PostInsert checks if the passed in item has a PostInsert method and invokes it.
func PostInsert(db DBHandle, item interface{}) error {
	if chk, ok := item.(PostInserter); ok {
		return chk.PostInsert(db)
	}
	return nil
}

// PreUpdate checks if the passed in item has a PreUpdate method and invokes it.
func PreUpdate(db DBHandle, item interface{}) error {
	if chk, ok := item.(PreUpdater); ok {
		return chk.PreUpdate(db)
	}
	return nil
}

// PostUpdate checks if the passed in item has a PostUpdate method and invokes it.
func PostUpdate(db DBHandle, item interface{}) error {
	if chk, ok := item.(PostUpdater); ok {
		return chk.PostUpdate(db)
	}
	return nil
}

// PreDelete checks if the passed in item has a PreDelete method and invokes it.
func PreDelete(db DBHandle, item interface{}) error {
	if chk, ok := item.(PreDeleter); ok {
		return chk.PreDelete(db)
	}
	return nil
}

// PostDelete checks if the passed in item has a PostDelete method and invokes it.
func PostDelete(db DBHandle, item interface{}) error {
	if chk, ok := item.(PostDeleter); ok {
		return chk.PostDelete(db)
	}
	return nil
}
