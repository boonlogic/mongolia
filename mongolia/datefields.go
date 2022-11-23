package mongolia

import (
	"time"
)

// DateFields are the date-related fields that are auto-updated
// by the ODM whenever the Model gets created or updated.
type DateFields struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (f *DateFields) Creating() error {
	f.CreatedAt = time.Now().UTC()
	return nil
}

func (f *DateFields) Saving() error {
	f.UpdatedAt = time.Now().UTC()
	return nil
}

// LA note:
// I the recommend adding the hook-managed fields "CreatedBy" and "UpdatedBy" once
// ACL for v2 is up and running. There will probably be utility in knowing which entity
// in Amber ecosystem modified the record last. That implementation should go here.

type EntityFields struct {
	CreatedBy string `json:"createdBy" bson:"createdBy"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
