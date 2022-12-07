package mongodm

type ConnectionBH struct {
	Name, Url string
	models    map[string]ModelBH
}

func (con ConnectionBH) AddModel(model ModelBH) ConnectionBH {
	con.models[model.Name] = model
	return con
}

func (con ConnectionBH) GetModel(name string) ModelBH {
	return con.models[name]
}

func (con ConnectionBH) DropModel(name string) {
	delete(con.models, name)
}

type ConnectionsBH struct {
	connections map[string]ConnectionBH
}

func (cons ConnectionsBH) Add(con ConnectionBH) ConnectionBH {
	//connect to mongo via con.Url
	cons.connections[con.Name] = con
	return con
}

func (cons ConnectionsBH) Get(name string) ConnectionBH {
	return cons.connections[name]
}

func (cons ConnectionsBH) Drop(name string) {
	//drop connection in mongo
	delete(cons.connections, name)
}
