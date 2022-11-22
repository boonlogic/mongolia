package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func validateIndexes(have []Index, want []Index) error {
	for _, wanted := range want {
		for _, idx := range have {
			if idx.Name == wanted.Name && !idx.Equals(wanted) {
				return errors.New(fmt.Sprintf("an index named '%s' exists, but does not match the desired specifications", idx.Name))
			}
		}
	}
	return nil
}

func missingIndexes(have []Index, want []Index) []Index {
	missing := make([]Index, 0)
	for _, wanted := range want {
		for _, idx := range have {
			if idx.Equals(wanted) {
				continue
			}
		}
		missing = append(missing, wanted)
	}
	return missing
}

func addIndexSet(coll *mongo.Collection, set []Index) error {
	var added = make([]string, 0) // keep track of indexes added for cleanup in case of error
	for _, idx := range set {
		err := addIndex(coll, idx)
		if err != nil {
			msg := fmt.Sprintf("failed to add index '%s': %s", idx.Name, err)
			for i, name := range added {
				if err := dropIndex(coll, name); err != nil {
					msg += fmt.Sprintf("\n\nwhile rolling back created indexes, another error occurred: %s", err)
					msg += fmt.Sprintf("\n\nleft behind %d dangling indexes", len(name)-i)
					panic(msg)
				}
			}
			return errors.New(msg)
		}
		added = append(added, idx.Name)
	}
	return nil
}

func prepareIndexes(coll *mongo.Collection, schema *Schema) ([]Index, error) {
	required := requiredIndexes(schema)
	existing, err := listIndexes(coll)
	if err != nil {
		return nil, err
	}
	if err = validateIndexes(existing, required); err != nil {
		return nil, err
	}
	missing := missingIndexes(existing, required)
	if err = addIndexSet(coll, missing); err != nil {
		return nil, err
	}
	all := append(existing, missing...)
	return all, nil
}
