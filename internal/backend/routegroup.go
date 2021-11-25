package backend

import (
	"log"
	"sanctuary/internal/db"
	"sanctuary/internal/form"
	"sanctuary/internal/utils"
)

type RouteGroupHandler struct{}

func NewRouteGroupHandler() *RouteGroupHandler {
	return &RouteGroupHandler{}
}

func (r *RouteGroupHandler) List(ctx Context) error {
	rgs, err := db.RouteGroups.Get(ctx.Request().Context())
	if err != nil {
		return ctx.ServerError()
	}
	return ctx.Success(rgs)
}

func (r *RouteGroupHandler) New(ctx Context, f form.NewRouteGroup) error {
	log.Println("new")
	if err := utils.CheckURL(f.Path); err != nil {
		return ctx.Error(400, err.Error())
	}
	if err := db.RouteGroups.Create(ctx.Request().Context(), &db.NewRouteGroupOption{
		Name:   f.Name,
		Path:   f.Path,
		Plugin: f.Plugin,
	}); err != nil {
		if err == db.ErrRouteGroupAlreadyExists {
			return ctx.Error(400, err.Error())
		} else {
			return ctx.ServerError()
		}
	}
	return nil
}
