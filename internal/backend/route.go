package backend

import (
	"log"
	"sanctuary/internal/db"
	"sanctuary/internal/form"
)

type RouteHandler struct{}

func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

func (r *RouteHandler) List(ctx Context, f form.GetRouteOption) error {
	routes, err := db.Routes.Get(ctx.Request().Context(), &db.GetRouteOption{
		OrderBy:  f.OrderBy,
		Page:     f.Page,
		PageSize: f.PageSize,
	})
	if err != nil {
		log.Println("List Routes error")
		return ctx.ServerError()
	}
	return ctx.Success(routes)
}

func (r *RouteHandler) New(ctx Context, f form.NewRoute) error {
	err := db.Routes.Create(ctx.Request().Context(), &db.NewRouteOption{
		Method:    f.Method,
		Path:      f.Path,
		GroupID:   f.GroupID,
		ServiceID: f.ServiceID,
	})
	switch err {
	case nil:
		return ctx.Success()
	case db.ErrRouteGroupNotExists:
		return ctx.Error(400, "Route group not exists.")
	case db.ErrRouteAlreadyExists:
		return ctx.Error(400, "Route path already exists.")
	case db.ErrNotHTTPMethod:
		return ctx.Error(400, "Route is not the HTTP method.")
	default:
		log.Println("RouteHandler.New: unknown error")
		return ctx.ServerError()
	}
}
