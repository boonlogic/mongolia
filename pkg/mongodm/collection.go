package mongodm

// Singleton containing all collections.
var collections = make(map[string]Collection)

type Collection interface {
	// CreateOne adds a document to the collection.
	CreateOne(Attributes) (*Document, error)

	// CreateMany adds multiple documents to the collection.
	CreateMany([]Attributes) ([]Document, error)

	// FindOne finds a document.
	FindOne(Query) (*Document, error)

	// FindMany finds multiple documents.
	FindMany(Query) ([]Document, error)

	// UpdateOne updates a document.
	UpdateOne(Query, Attributes) (*Document, error)

	// UpdateMany updates multiple documents.
	UpdateMany(Query, Attributes) ([]Document, error)

	// RemoveOne deletes a document.
	RemoveOne(Query) (*Document, error)

	// RemoveMany deletes many documents.
	RemoveMany(Query) ([]Document, error)
}

func GetCollection(name string) (Collection, bool) {
	c, ok := collections[name]
	if !ok {
		return nil, false
	}
	return c, true
}
