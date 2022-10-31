package mgm

var odm = new(ODM)

type ODM struct {
}

func (o *ODM) Collection(*Model) *Collection {
	return nil
}
