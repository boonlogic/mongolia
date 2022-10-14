package odm

type ModelRegistry map[string]Model

var modelRegistry = ModelRegistry{}

func GetModel(name string) Model {
	return modelRegistry[name]
}

type Model interface {
	preValidate(any) *Model
	postValidate(any) *Model
	preSave(any) *Model
	postSave(any) *Model
	preCreate(any) *Model
	postCreate(any) *Model
	preUpdate(any) *Model
	postUpdate(any) *Model
	preRemove(any) *Model
	postRemove(any) *Model

	CreateOne(any) (*Document, error)
	CreateMany(any) ([]Document, error)
	FindOne(any) (*Document, error)
	FindMany(any) ([]Document, error)
	UpdateOne(any) (*Document, error)
	UpdateMany(any) ([]Document, error)
	RemoveOne(any) (*Document, error)
	RemoveMany(any) ([]Document, error)
}

func RegisterModel(name string, spec []byte, hooks *Hooks) error {
	vfunc, err := validateSpec(name, spec)
	if err != nil {
		return err
	}
	s := Schema{
		Name:       name,
		Definition: spec,
		Validator:  vfunc,
		Hooks:      hooks,
	}
	modelRegistry[name] = s
	return nil
}
