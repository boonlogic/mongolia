package v0

// BaseModel implements the Model interface.
// It is embedded at the top level of a custom struct to make that struct a Model.
type BaseModel struct {
	Document
	coll *Collection
}

//func (m BaseModel) Save(opts *options2.SaveOptions) error {
//	filter := bson.M{
//		"id": m.Document["id"],
//	}
//	m.coll.updateOne(filter, m)
//	return nil
//}
//
//func (m BaseModel) Remove(opts *options2.RemoveOptions) error {
//	filter := bson.M{
//		"id": m.Document["id"],
//	}
//	m.coll.removeOne(filter, m)
//	return nil
//}

// todo
//func (d BaseModel) Populate(opts *options.PopulateOptions) error {
//	return nil
//}

func (m BaseModel) Get(key string) any {
	return m.Document[key]
}

func (m BaseModel) Set(key string, val any) {
	m.Document[key] = val
}