package main

import (
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"testing"
)

//type Tenant struct {
//	mongolia.DefaultModel `bson:",inline"`
//	TenantID              string `json:"tenantId" bson:"tenantId"`
//	Name                  string `json:"name" bson:"name"`
//}
//
//func NewTenant(name string) *Tenant {
//	return &Tenant{
//		Name: name,
//	}
//}

func TestSmoke(t *testing.T) {
	cfg := mongolia.NewConfig().
		SetURI("mongodb://localhost:27017").
		SetDBName("mongolia-local")

	err := mongolia.Connect(cfg)
	require.Nil(t, err)

	err = mongolia.AddSchema("tenant", "test/tenant.json")
	require.Nil(t, err)

	coll, err := mongolia.GetCollection("tenant")
	require.NotNil(t, coll)
	require.Nil(t, err)
}
