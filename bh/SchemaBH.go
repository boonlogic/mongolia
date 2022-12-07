package mongodm

func crudVal(attributes map[string]any, spec SpecBH) bool {
	spec = spec + "spec"
	attributes["abc"] = 123
	return true
}

type SchemaBH struct {
	ReadSpec, CreateSpec, UpdateSpec SpecBH
}

func (s SchemaBH) ValidateRead(attributes map[string]any) bool {
	//validate against ReadSpec
	return crudVal(attributes, "read")
}

func (s SchemaBH) ValidateCreate(attributes map[string]any) bool {
	//validate against FindSpec
	return crudVal(attributes, "create")
}

func (s SchemaBH) ValidateUpdate(attributes map[string]any) bool {
	//validate against UpdateSpec
	return crudVal(attributes, "update")
}

func (s SchemaBH) PreValidate(attributes map[string]any) map[string]any {
	return attributes
}

func (s SchemaBH) PostValidate(attributes map[string]any) map[string]any {
	return attributes
}

func (s SchemaBH) PreCreate(attributes map[string]any) map[string]any {
	return attributes
}

func (s SchemaBH) PostCreate(document DocumentBH) DocumentBH {
	return document
}

func (s SchemaBH) PreUpdate(attributes map[string]any) map[string]any {
	return attributes
}

func (s SchemaBH) PostUpdate(document DocumentBH) DocumentBH {
	return document
}

func (s SchemaBH) PreSave(attributes map[string]any) map[string]any {
	return attributes
}

func (s SchemaBH) PostSave(document DocumentBH) DocumentBH {
	return document
}

func (s SchemaBH) PreDelete(query QueryBH) bool {
	query.Exec()
	return true
}

func (s SchemaBH) PostDelete(document DocumentBH) DocumentBH {
	return document
}

// ------------ Spec is a block of OpenAPI 3 -----------

type SpecBH string
