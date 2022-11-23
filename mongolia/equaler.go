package mongolia

// Any struct with this Equals method satisfies the Equaler interface.
type Equaler interface {
	Equals(other Equaler) bool
}

// Equal defines the equality relation for Model. The relation is commutative.
// This makes it so users of Model do not need to worry about comparison order.
func Equal(a Model, b Model) bool {
	return a.Equals(b) && b.Equals(a)
}
