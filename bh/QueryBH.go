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

func (q QueryBH) Exec() {
	//meant to be protected. not public or private. Mongolia code only can call it
	//do all logic here
	//build up compound mongo driver
}
