package backend

import (
	"sanctuary/internal/db"
	"sanctuary/internal/form"
)

type ServiceGroupHandler struct{}

func NewServiceGroupHandler() *ServiceGroupHandler {
	return &ServiceGroupHandler{}
}

func (r *ServiceGroupHandler) New(ctx Context, f form.NewServiceGroup) error {
	if err := db.ServiceGroups.NewServiceGroup(&db.NewServiceGroupOption{
		Name:   f.Name,
		Plugin: f.Plugin,
	}); err != nil {
		if err == db.ErrServiceGroupAlreadyExists {
			return ctx.Error(400, err.Error())
		} else {
			return ctx.ServerError()
		}
	}
	return nil
}
