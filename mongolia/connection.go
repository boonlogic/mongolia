package mongolia

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(config *Config) error {
	opts := options.Client().ApplyURI(*config.URI)
	mgmconf := &mgm.Config{
		CtxTimeout: *config.Timeout,
	}
	if err := mgm.SetDefaultConfig(mgmconf, *config.DBName, opts); err != nil {
		return err
	}
	return nil
}

func connectCollection(name string, schema *Schema) (*Collection, error) {
	coll := mgm.CollectionByName(name)

	_ = schema.RequiredIndexes()
	_, err := listIndexes(coll.Collection)
	if err != nil {
		return nil, err
	}

	// todo:
	// [Add indexes that are needed to match the schema].
	//
	// [Get indexes required by the schema].
	// required = RequiredIndexes()

	// [Cross off any that are already in the collection].
	// exists_i []
	// for i, collidx in listIndexes(coll)
	//     for j, reqidx in required:
	//         if reqidx == collidx:
	//             exists_i.append(i)
	// del required[exists_i]
	//
	// [Iterate through remaining required indexes and add them].
	// created = []
	// for i, idx in required:
	//     try:
	//         addIndex(coll, idx)
	//     except:
	//         [If an error occurs, delete those that were previously added].
	//         for name in created:
	//             try:
	//                 dropIndex(coll, name)
	//             except:
	//                 panic("failed to rollback previously created indexes while handling index creation failure")
	//     added_names.append(idx.name)

	c := &Collection{
		schema: schema,
		coll:   coll,
	}
	return c, nil
}
