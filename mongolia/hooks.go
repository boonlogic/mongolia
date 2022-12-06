package mongolia

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// PreCreateHook is called before saving a new model to the database
type PreCreateHook interface {
	PreCreate(context.Context) error
}

// CreatedHook is called after a model has been created
type PostCreateHook interface {
	PostCreate(context.Context) error
}

// PreUpdateHook is called before updating a model
type PreUpdateHook interface {
	PreUpdate(context.Context) error
}

// PostUpdateHook is called after a model is updated
type PostUpdateHook interface {
	PostUpdate(ctx context.Context, result *mongo.UpdateResult) error
}

// PreSaveHook is called before a model (new or existing) is saved to the database.
type PreSaveHook interface {
	PreSave(context.Context) error
}

// PostSaveHook is called after a model is saved to the database.
type PostSaveHook interface {
	PostSave(context.Context) error
}

// PreDeleteHook is called before a model is deleted
type PreDeleteHook interface {
	PreDelete(context.Context) error
}

// PostDeleteHook is called after a model is deleted
type PostDeleteHook interface {
	PostDelete(ctx context.Context, result *mongo.DeleteResult) error
}

func beforeCreateHooks(ctx context.Context, model Model) error {
	if hook, ok := model.(PreCreateHook); ok {
		if err := hook.PreCreate(ctx); err != nil {
			return err
		}
	}
	if hook, ok := model.(PreSaveHook); ok {
		if err := hook.PreSave(ctx); err != nil {
			return err
		}
	}

	return nil
}

func afterCreateHooks(ctx context.Context, model Model) error {
	if hook, ok := model.(PostCreateHook); ok {
		if err := hook.PostCreate(ctx); err != nil {
			return err
		}
	}

	if hook, ok := model.(PostSaveHook); ok {
		if err := hook.PostSave(ctx); err != nil {
			return err
		}
	}

	return nil
}

func beforeUpdateHooks(ctx context.Context, model Model) error {
	if hook, ok := model.(PreUpdateHook); ok {
		if err := hook.PreUpdate(ctx); err != nil {
			return err
		}
	}

	if hook, ok := model.(PreSaveHook); ok {
		if err := hook.PreSave(ctx); err != nil {
			return err
		}
	}

	return nil
}

func afterUpdateHooks(ctx context.Context, updateResult *mongo.UpdateResult, model Model) error {
	if hook, ok := model.(PostUpdateHook); ok {
		if err := hook.PostUpdate(ctx, updateResult); err != nil {
			return err
		}
	}

	if hook, ok := model.(PostSaveHook); ok {
		if err := hook.PostSave(ctx); err != nil {
			return err
		}
	}

	return nil
}

func beforeDeleteHooks(ctx context.Context, model Model) error {
	if hook, ok := model.(PreDeleteHook); ok {
		if err := hook.PreDelete(ctx); err != nil {
			return err
		}
	}

	return nil
}

func afterDeleteHooks(ctx context.Context, deleteResult *mongo.DeleteResult, model Model) error {
	if hook, ok := model.(PostDeleteHook); ok {
		if err := hook.PostDelete(ctx, deleteResult); err != nil {
			return err
		}
	}

	return nil
}
