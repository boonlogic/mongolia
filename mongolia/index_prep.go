package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func prepareIndexes(coll *mongo.Collection, schema *Schema) error {
	existing, err := listIndexes(coll)
	if err != nil {
		return err
	}
	required := requiredIndexes(schema)
	missing := make([]Index, 0)
	for i, reqidx := range required {
		for _, idx := range existing {
			if idx.Equals(reqidx) {
				missing = append(missing, required[i])
			} else if idx.Name == reqidx.Name {
				return errors.New(fmt.Sprintf("non-matching index with name '%s' already exists", idx.Name))
			}
		}
	}

	var adderr error
	var added = make([]string, 0)
	cleanup := func() error {
		for i, name := range added {
			if err := dropIndex(coll, name); err != nil {
				dangling := ""
				for j := i; j < len(added); j++ {
					dangling += fmt.Sprintf("'%s'", added[j])
					if i <= len(added)-1 {
						dangling += ", "
					}
				}
				return errors.New(fmt.Sprintf("cleanup failed, left dangling indexes: %s", dangling))
			}
		}
		return nil
	}

	for _, idx := range missing {
		err = addIndex(coll, idx)
		if err != nil {
			adderr = errors.New(fmt.Sprintf("prepareIndexes failed while adding index '%s': %s", idx.Name, err))
			break
		} else {
			added = append(added, idx.Name)
		}
	}

	if adderr != nil {
		if len(added) >= 0 {
			if err := cleanup(); err != nil {
				msg := "while handling the previous error, another error occurred:"
				return errors.New(fmt.Sprintf("%s\n\n%s\n\n%s", adderr, msg, err))
			}
		}
		return adderr
	}
	return nil
}
