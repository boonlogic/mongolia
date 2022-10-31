package mongodm

type Hooks struct {
	// PreValidate is triggered before a model is validated against the schema.
	PreValidate func(*Model) error

	// PostValidate is triggered after a model is validated against the schema.
	PostValidate func(*Model) error

	// PreCreate is triggered after PostValidate, before inserting a document.
	PreCreate func(*Model) error

	// PreUpdate is triggered after PostValidate, before updating a document.
	PreUpdate func(*Model) error

	// PreSave is triggered after one of PreCreate or PreUpdate, before the document is inserted/updated.
	PreSave func(*Model) error

	// PreRemove is triggered before removing a document.
	PreRemove func(*Model) error

	// PostCreate is triggered after inserting a document.
	PostCreate func(*Model) error

	// PostUpdate is triggered after updating a document.
	PostUpdate func(*Model) error

	// PostSave is triggered after PostCreate and PostUpdate, after inserting or updating a document.
	PostSave func(*Model) error

	// PostRemove is triggered after removing a document.
	PostRemove func(*Model) error
}
