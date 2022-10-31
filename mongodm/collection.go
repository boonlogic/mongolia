package mongodm

import (
	"context"
	"github.com/Kamva/mgm"
)

type Collection struct {
	*mgm.Collection
}

func (c *Collection) Create(m *Model) error {
	return nil
}

func (c *Collection) Update(m *Model) error {
	return nil
}

func (c *Collection) Delete(m *Model) error {
	return nil
}

func (c *Collection) CreateWithCtx(ctx context.Context, model *Model) error {
	var m mgm.Model
	m = *model

	mgmModel := mgm.Model(*model)

	return c.CreateWithCtx(m * Model)
	return nil
}

func (c *Collection) UpdateWithCtx(m *Model) error {
	return nil
}

func (c *Collection) DeleteWithCtx(m *Model) error {
	return nil
}
