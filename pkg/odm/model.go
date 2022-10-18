package odm

type ModelRegistry map[string]Model

var modelRegistry = ModelRegistry{}

func GetModel(name string) Model {
	return modelRegistry[name]
}

type Model interface {
	// preValidate is triggered before a document is validated against the schema.
	preValidate() error

	// postValidate is triggered after a document is validate against the schema.
	postValidate() error

	// preCreate is triggered after postValidate and before inserting a document.
	preCreate() error

	// preUpdate is triggered after postValidate and before updating a document.
	preUpdate() error

	// preSave is triggered after preCreate/preUpdate and before inserting or updating a document.
	preSave() error

	// preRemove is triggered before removing a document.
	preRemove() error

	// postCreate is triggered after inserting a document.
	postCreate() error

	// postUpdate is triggered after updating a document.
	postUpdate() error

	// postSave is triggered after postCreate/postUpdate, after inserting or updating a document.
	postSave() error

	// postRemove is triggered after removing a document.
	postRemove() error

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
