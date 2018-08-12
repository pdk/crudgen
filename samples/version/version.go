package version

import (
	"time"

	"github.com/google/uuid"
)

// Version models have basic fields of uuid, version_id and version_at.
type Version struct {
	UUID          uuid.UUID // the object id
	VersionID     int64     `crud:"autoincr"`         // which version (actual sequence/row id)
	VersionAt     time.Time `crud:"update_timestamp"` // when this version was created
	ActiveVersion bool      // flag indicating which is the currently active version
}
