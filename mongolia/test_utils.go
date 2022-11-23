package mongolia

import "go.mongodb.org/mongo-driver/bson/primitive"

func Drop() {
	odm.drop()
}

func NewOID() string {
	return primitive.NewObjectID().Hex()
}

type Tenant struct {
	DefaultModel `bson:",inline"`
	TenantID     *string `json:"tenantId,omitempty" bson:"tenantId,omitempty"`
	Name         *string `json:"name,omitempty" bson:"name,omitempty"`
}

func NewTenant(id string, name string) *Tenant {
	return &Tenant{
		TenantID: &id,
		Name:     &name,
	}
}

func (t *Tenant) Equals(other *Tenant) bool {
	if *t.TenantID != *other.TenantID {
		return false
	}
	if *t.Name != *other.Name {
		return false
	}
	return true
}
