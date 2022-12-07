package mongolia

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Query struct {
	sort     bson.D
	fields   bson.D
	limit    int
	skip     int
	populate string
	where    bson.D
	count    bool
}

func (q Query) Where(att bson.D) Query {
	q.where = att
	return q
}

func (q Query) Sort(att bson.D) Query {
	q.sort = att
	return q
}

func (q Query) Populate(att string) Query {
	q.populate = att
	return q
}

func (q Query) Select(fields bson.D) Query {
	q.fields = fields
	return q
}

func (q Query) Limit(limit int) Query {
	q.limit = limit
	return q
}

func (q Query) Skip(skip int) Query {
	q.skip = skip
	return q
}

func (q Query) pipeline(model Model) mongo.Pipeline {
	//handle populate
	tagReferences := model.GetTagReferences()
	fromcollection, ok := tagReferences[q.populate]
	if !ok {
		fromcollection = q.populate
	}
	foreign := "_id"

	matchStage := bson.D{{"$match", q.Where}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", fromcollection}, {"localField", q.populate}, {"foreignField", foreign}, {"as", q.populate}}}}
	projectStage := bson.D{{"$project", q.Select}}
	limitStage := bson.D{{"$limit", q.skip + q.limit}}
	skipStage := bson.D{{"$skip", q.skip}}
	sortStage := bson.D{{"$sort", q.sort}}

	return mongo.Pipeline{matchStage, lookupStage, projectStage, limitStage, skipStage, sortStage}
}
