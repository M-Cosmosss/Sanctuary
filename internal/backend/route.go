package backend

import (
	"log"
	"path"
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
	g, err := db.RouteGroups.GetByID(ctx.Request().Context(), f.GroupID)
	switch err {
	case nil:
		break
	case db.ErrRouteGroupNotExists:
		return ctx.Error(404, "RouteGroup ID not found.")
	default:
		return ctx.ServerError()
	}
	err = db.Routes.Create(ctx.Request().Context(), &db.NewRouteOption{
		Method:    f.Method,
		FullPath:  path.Join(g.Path, f.Path),
		Path:      f.Path,
		GroupID:   g.ID,
		ServiceID: f.ServiceID,
	})
	switch err {
	case nil:
		return ctx.Success()
	case db.ErrRouteAlreadyExists:
		return ctx.Error(400, "Route path already exists.")
	default:
		return ctx.ServerError()
	}
}
