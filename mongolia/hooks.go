package mongolia

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func beforeCreateHooks(model Model) error {
	if err := model.PreCreate(); err != nil {
		return err
	}
	if err := model.PreSave(); err != nil {
		return err
	}

	return nil
}

func afterCreateHooks(model Model) error {
	if err := model.PostCreate(); err != nil {
		return err
	}
	if err := model.PostSave(); err != nil {
		return err
	}
	return nil
}

func beforeUpdateHooks(model Model) error {
	if err := model.PreUpdate(); err != nil {
		return err
	}
	if err := model.PreSave(); err != nil {
		return err
	}

	return nil
}

func afterUpdateHooks(updateResult *mongo.UpdateResult, model Model) error {
	if err := model.PostUpdate(updateResult); err != nil {
		return err
	}
	if err := model.PostSave(); err != nil {
		return err
	}

	return nil
}

func beforeDeleteHooks(model Model) error {
	if err := model.PreDelete(); err != nil {
		return err
	}

	return nil
}

func afterDeleteHooks(deleteResult *mongo.DeleteResult, model Model) error {
	if err := model.PostDelete(deleteResult); err != nil {
		return err
	}

	return nil
}
