package mongolia

// Equaler tells whether one Model is equal to another
type Equaler interface {
	Equals(other Model) bool
}
