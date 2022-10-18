package odm

type Hooks struct {
	PreValidate  func() error
	PostValidate func() error
	PreSave      func() error
	PostSave     func() error
	PreCreate    func() error
	PostCreate   func() error
	PreUpdate    func() error
	PostUpdate   func() error
	PreRemove    func() error
	PostRemove   func() error
}
