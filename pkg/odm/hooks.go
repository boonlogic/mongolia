package odm

type Hooks struct {
	PreValidate  func(any) *Model
	PostValidate func(any) *Model
	PreSave      func(any) *Model
	PostSave     func(any) *Model
	PreCreate    func(any) *Model
	PostCreate   func(any) *Model
	PreUpdate    func(any) *Model
	PostUpdate   func(any) *Model
	PreRemove    func(any) *Model
	PostRemove   func(any) *Model
}

