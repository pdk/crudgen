package crudlib

// PreInserter offers a pre-insert operation which might return an error to
// indicate the operation should be aborted.
type PreInserter interface {
	PreInsert() error
}

// PreUpdater offers a pre-update operation which might return an error to
// indicate the operation should be aborted.
type PreUpdater interface {
	PreUpdate() error
}

// PostDeleter offers a post-deletion operation.
type PostDeleter interface {
	PostDelete()
}

// PreInsert checks if the passed in item has a PreInsert method and invokes it.
func PreInsert(item interface{}) error {
	if chk, ok := item.(PreInserter); ok {
		return chk.PreInsert()
	}

	return nil
}

// PreUpdate checks if the passed in item has a PreUpdate method and invokes it.
func PreUpdate(item interface{}) error {
	if chk, ok := item.(PreUpdater); ok {
		return chk.PreUpdate()
	}

	return nil
}

// PostDelete checks if the passed in item has a PostDelete method and invokes it.
func PostDelete(item interface{}) {
	if chk, ok := item.(PostDeleter); ok {
		chk.PostDelete()
	}
}
