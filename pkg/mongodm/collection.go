package mongodm

// Singleton containing all collections.
var collections = make(map[string]Collection)

type Collection interface {
	// CreateOne adds a document to the collection.
	CreateOne(any) (*Document, error)

	// CreateMany adds multiple documents to the collection.
	CreateMany(any) ([]Document, error)

	// FindOne finds a document.
	FindOne(any) (*Document, error)

	// FindMany finds multiple documents.
	FindMany(any) ([]Document, error)

	// UpdateOne updates a document.
	UpdateOne(any) (*Document, error)

	// UpdateMany updates multiple documents.
	UpdateMany(any) ([]Document, error)

	// RemoveOne deletes a document.
	RemoveOne(any) (*Document, error)

	// RemoveMany deletes many documents.
	RemoveMany(any) ([]Document, error)
}

func GetCollection(name string) (Collection, bool) {
	s, ok := collections[name]
	if !ok {
		return nil, false
	}
	m := Collection(s)
	return m, true
}
