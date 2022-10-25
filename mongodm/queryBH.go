package mongodm

import (
	"go.mongodb.org/mongo-driver/bson"
)

type QueryBH struct {
	sort, fields []string
	limit, skip  int
	populate     []Populate
	where        []bson.M
	count        bool
}

type Populate struct {
	Path   string
	Select string
	Where  QueryBH
}

func (q QueryBH) Where(att bson.M) QueryBH {
	q.where = append(q.where, att)
	return q
}

func (q QueryBH) Sort(att string) QueryBH {
	q.sort = append(q.sort, att)
	return q
}

func (q QueryBH) Populate(att Populate) QueryBH {
	q.populate = append(q.populate, att)
	return q
}

func (q QueryBH) Select(attributes string) QueryBH {
	q.fields = append(q.fields, attributes)
	return q
}

func (q QueryBH) Limit(limit int) QueryBH {
	q.limit = limit
	return q
}

func (q QueryBH) Skip(skip int) QueryBH {
	q.skip = skip
	return q
}

func (q QueryBH) Count(count bool) QueryBH {
	q.count = count
	return q
}

func (q QueryBH) exec() {
	//meant to be protected. not public or private. Mongolia code only can call it
	//do all logic here
	//build up compound mongo driver
}

///===============================

func Samples() {

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
	myQuery.exec()
	sameQuery.exec()

}
