package mongodm

type Hooks struct {
	PreValidate  func(*Document) error
	PostValidate func(*Document) error
	PreSave      func(*Document) error
	PostSave     func(*Document) error
	PreCreate    func(*Document) error
	PostCreate   func(*Document) error
	PreUpdate    func(*Document) error
	PostUpdate   func(*Document) error
	PreRemove    func(*Document) error
	PostRemove   func(*Document) error
}