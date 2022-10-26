package mongodm

type Index []IndexField

type IndexField struct {
	Key        string
	Increasing bool
}
