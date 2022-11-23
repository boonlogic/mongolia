package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	cfg := mongolia.NewConfig().
		SetURI("mongodb://localhost:27017").
		SetDBName("mongolia-local-tmp").
		SetTimeout(10 * time.Second).
		SetEphemeral(true)

	err := mongolia.Connect(cfg)
	require.Nil(t, err)

	mongolia.Drop()

	// try to get a nonexistent collection
	coll, err := mongolia.GetCollection("nonexistent")
	require.NotNil(t, err)
	require.Nil(t, coll)

	// add a schema
	path := os.Getenv("SCHEMA_PATH")
	if path == "" {
		path = "test/tenant.json"
	}
	err = mongolia.AddSchema("tenant", "test/tenant.json")
	require.Nil(t, err)

	// get the corresponding collection
	coll, err = mongolia.GetCollection("tenant")
	require.Nil(t, err)
	require.NotNil(t, coll)

	// try to find a nonexistent tenant
	var result *mongolia.Tenant
	q := map[string]any{
		"tenantId": mongolia.NewOID(),
	}
	coll.First(q, result, nil)
	require.Nil(t, err)
	require.Nil(t, result)

	// make a new tenant in memory
	var tenant *mongolia.Tenant
	tenantId := mongolia.NewOID()
	name := "luke"
	tenant = mongolia.NewTenant(tenantId, name)

	// create that tenant in the database
	err = coll.Create(tenant, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	require.Nil(t, err)
	require.NotNil(t, tenant)
	require.Equal(t, tenantId, *tenant.TenantID)
	require.Equal(t, name, *tenant.Name)

	// find the created tenant
	var found = &mongolia.Tenant{}
	q = map[string]any{
		"tenantId": tenantId,
	}
	err = coll.First(q, found, nil)
	require.Nil(t, err)
	require.True(t, found.Equals(tenant))

	//// change the tenant in memory does not change it in the database
	//tid2 := NewOID()
	//name2 := "brad"
	//tenant.TenantID = tid2
	//tenant.Name = name2

}
