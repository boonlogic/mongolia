package mongodm

func validate(spec SpecBH) bool {
	spec = spec + "spec"
	return true
}

type ModelBH struct {
	Name   string
	Schema SchemaBH
}

func (bm ModelBH) Create(attributes map[string]any) DocumentBH {
	//run pre validate
	attributes = bm.Schema.PreValidate(attributes)
	//validate via schemas
	bm.Schema.ValidateCreate(attributes)
	//run post validate hook
	attributes = bm.Schema.PostValidate(attributes)
	//run pre create hook
	attributes = bm.Schema.PreCreate(attributes)
	//create record in mongo
	ats := make(map[string]any)
	ats["abc"] = 123
	doc := DocumentBH{bm, attributes}
	//run post create hook
	bm.Schema.PostCreate(doc)
	return doc
}

func (bm ModelBH) Read(query QueryBH) []DocumentBH {
	//query mongo
	query.Exec()
	ats := make(map[string]any)
	ats["abc"] = 123
	doc := DocumentBH{bm, ats}
	docs := []DocumentBH{doc}
	return docs
}

func (bm ModelBH) Update(query QueryBH, attributes map[string]any) {
	//run before hooks
	attributes = bm.Schema.PreValidate(attributes)
	bm.Schema.ValidateUpdate(attributes)
	attributes = bm.Schema.PostValidate(attributes)
	attributes = bm.Schema.PreCreate(attributes)
	attributes = bm.Schema.PreUpdate(attributes)
	//update  mongo
	query.Exec()
	//update mongo
	attributes["abc"] = 123
	doc := DocumentBH{bm, attributes}
	//run schema post hooks
	bm.Schema.PostUpdate(doc)
	bm.Schema.PostCreate(doc)
}

func (bm ModelBH) Delete(query QueryBH) {
	//run schema pre hook
	bm.Schema.PreDelete(query)
	//delete from mongo
	query.Exec()
	ats := make(map[string]any)
	ats["abc"] = 123
	doc := DocumentBH{bm, ats}
	//run schema post hook
	bm.Schema.PostDelete(doc)
}

type ModelsBH struct {
	models map[string]ModelBH
}

func (mods ModelsBH) Add(name string, mod ModelBH) ModelBH {
	//validate the actual openapi 3 blocks
	validate(mod.Schema.ReadSpec)
	validate(mod.Schema.CreateSpec)
	validate(mod.Schema.UpdateSpec)
	//add collection to mongo if setting up indexes ...
	mods.models[name] = mod
	return mod
}

func (mods ModelsBH) Get(name string) ModelBH {
	return mods.models[name]
}

func (mods ModelsBH) Drop(name string) {
	//drop collection from mongodb
	delete(mods.models, name)
}
