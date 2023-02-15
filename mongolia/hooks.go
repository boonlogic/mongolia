package mongolia

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func beforeSaveHooks(model Model) *Error {
	if err := model.PreSave(); err != nil {
		return NewError(406, err)
	}

	return nil
}

func afterSaveHooks(model Model) *Error {
	if err := model.PostSave(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func afterReadHooks(model Model) *Error {
	if err := model.PostRead(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func beforeCreateHooks(model Model) *Error {
	if err := model.PreCreate(); err != nil {
		return NewError(406, err)
	}
	if err := model.PreSave(); err != nil {
		return NewError(406, err)
	}

	return nil
}

func afterCreateHooks(model Model) *Error {
	if err := model.PostCreate(); err != nil {
		return NewError(406, err)
	}
	if err := model.PostSave(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func beforeUpdateHooks(update any, model Model) *Error {
	if update == nil {
		if err := model.PreUpdateModel(); err != nil {
			return NewError(406, err)
		}
		if err := model.PreSave(); err != nil {
			return NewError(406, err)
		}
	} else {
		if err := model.PreUpdate(update); err != nil {
			return NewError(406, err)
		}
	}

	return nil
}

func afterUpdateHooks(updateResult *mongo.UpdateResult, model Model) *Error {
	if err := model.PostUpdateModel(updateResult); err != nil {
		return NewError(406, err)
	}
	if err := model.PostSave(); err != nil {
		return NewError(406, err)
	}
	return nil
}

func beforeDeleteHooks(model Model) *Error {
	if err := model.PreDelete(); err != nil {
		return NewError(406, err)
	}

	return nil
}

func afterDeleteHooks(deleteResult *mongo.DeleteResult, model Model) *Error {
	if err := model.PostDelete(deleteResult); err != nil {
		return NewError(406, err)
	}

	return nil
}
