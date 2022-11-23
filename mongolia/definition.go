package mongolia

type Definition []byte

func (d *Definition) UniqueFields() []string {
	// todo: implement

	// Parse the unique fields from a Boon-defined x-attribute of the schema, "x-unique".
	// x-unique is a list of fields whose values must be unique within that schema's collection.

	return []string{}
}
