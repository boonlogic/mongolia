package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
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
	defer mongolia.Drop()

	// no collections exist at startup
	coll, err := mongolia.GetCollection("nonexistent")
	require.NotNil(t, err)
	require.Nil(t, coll)

	// when a schema is added, a corresponding collection is made
	coll, err = mongolia.AddSchema("tenant", "test/tenant.json")
	require.Nil(t, err)

	mongolia.Drop()

	// get the collection for confirmation
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

	// create an associated record in the "tenants" collection
	err = coll.Create(tenant, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	require.Nil(t, err)
	require.NotNil(t, tenant)
	require.Equal(t, tenantId, *tenant.TenantID)
	require.Equal(t, name, *tenant.Name)

	// find the created tenant by ID
	var found = &mongolia.Tenant{}
	q = map[string]any{
		"tenantId": tenantId,
	}
	err = coll.FindByID(tenant.GetID(), found)
	require.Nil(t, err)
	require.Equal(t, found.GetID(), tenant.GetID())

	// find the created tenant using First
	q = map[string]any{
		"tenantId": tenantId,
	}
	err = coll.First(q, found, nil)
	require.Nil(t, err)
	require.Equal(t, found.GetID(), tenant.GetID())

	// change the tenant in memory
	tid2 := mongolia.NewOID()
	name2 := "brad"
	tenant.TenantID = &tid2
	tenant.Name = &name2

	// its DB document should remain unchanged
	err = coll.FindByID(tenant.GetID(), found)
	require.Nil(t, err)
	require.NotEqual(t, *found.TenantID, tid2)
	require.NotEqual(t, *found.Name, name2)

	// update the DB document to match the struct
	err = coll.Update(tenant, nil)
	require.Nil(t, err)

	// validate that Update did not change the tenant struct
	require.Equal(t, *tenant.TenantID, tid2)
	require.Equal(t, *tenant.Name, name2)

	// validate that the DB document matches the tenant struct
	err = coll.FindByID(tenant.GetID(), found)
	require.Nil(t, err)
	require.Equal(t, *found.TenantID, *tenant.TenantID)
	require.Equal(t, *found.Name, *tenant.Name)

	// delete the tenant
	err = coll.Delete(tenant)
	require.Nil(t, err)

	// make sure they're gone
	found = new(mongolia.Tenant)
	err = coll.FindByID(tenant.GetID(), found)
	require.NotNil(t, err)
	require.Empty(t, found)
}
