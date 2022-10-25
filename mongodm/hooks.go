package mongodm

type Hooks struct {
	// PreValidate is triggered before a document is validated against the schema.
	PreValidate func(*Document) error

	// PostValidate is triggered after a document is validated against the schema.
	PostValidate func(*Document) error

	// PreCreate is triggered after PostValidate and before inserting a document.
	PreCreate func(*Document) error

	// PreUpdate is triggered after PostValidate and before updating a document.
	PreUpdate func(*Document) error

	// PreSave is triggered after PreCreate or PreUpdate and before inserting or updating a document.
	PreSave func(*Document) error

	// PreRemove is triggered before removing a document.
	PreRemove func(*Document) error

	// PostCreate is triggered after inserting a document.
	PostCreate func(*Document) error

	// PostUpdate is triggered after updating a document.
	PostUpdate func(*Document) error

	// PostSave is triggered after PostCreate and PostUpdate, after inserting or updating a document.
	PostSave func(*Document) error

	// PostRemove is triggered after removing a document.
	PostRemove func(*Document) error
}
