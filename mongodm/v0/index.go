package v0

type Index struct {
	Name string
	Keys []IndexKey
}

type IndexKey struct {
	Field      string
	Increasing bool
}
