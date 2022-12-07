package mongodm

import (
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test(t *testing.T) {

	var myQuery = QueryBH{}

	myQuery.
		Where(bson.M{"active": true, "confirmed": true}).
		Limit(25).
		Skip(25).
		Select("id username email name tenancies").
		Populate(Populate{Path: "tenancies.tenant", Select: "id name"}).
		Populate(Populate{Path: "tenancies.role.perms", Select: "key name"}).
		Sort("-createdAt name")

	//or like this
	var sameQuery = QueryBH{}
	sameQuery.
		Where(bson.M{"active": true}).
		Where(bson.M{"confirmed": true}).
		Limit(25).
		Skip(25).
		Select("id").
		Select("username").
		Select("email").
		Select("name tenancies").
		Select("tenancies").
		Populate(Populate{Path: "tenancies.tenant", Select: "id name"}).
		Populate(Populate{Path: "tenancies.role"}).
		Populate(Populate{Path: "tenancies.role.perms", Select: "key name"}).
		Sort("-createdAt").
		Sort("name")

	// usage as params in find var recs1 documents[] = collectionXYZ.findMany(myQuery)
	// usage as params in find var recs2 documents[] = collectionXYZ.findMany(sameQuery)
	myQuery.Exec()
	sameQuery.Exec()

}
