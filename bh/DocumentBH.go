package mongodm

import "go.mongodb.org/mongo-driver/bson"

type PageBH struct {
	Skip  int
	Limit int
}

type CountBH struct {
	Filtered   int
	Collection int
}

type DocumentsBH struct {
	Documents []DocumentBH
	Count     CountBH
	Page      PageBH
}

type DocumentBH struct {
	Model      ModelBH
	Attributes map[string]any
}

func (s DocumentBH) Save() {
	q := QueryBH{}
	q.Where(bson.M{"id": s.Attributes["id"]})
	s.Model.Update(q, s.Attributes)
}

func (s DocumentBH) Delete() {
	q := QueryBH{}
	q.Where(bson.M{"id": s.Attributes["id"]})
	s.Model.Delete(q)
}
