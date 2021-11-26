package backend

import (
	"log"
	"sanctuary/internal/db"
	"sanctuary/internal/form"
)

type ServiceHandler struct{}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

func (s *ServiceHandler) List(ctx Context, f form.GetServiceOption) error {
	services, err := db.Services.Get(ctx.Request().Context(), &db.GetServiceOption{
		OrderBy:  f.OrderBy,
		Page:     f.Page,
		PageSize: f.PageSize,
	})
	if err != nil {
		log.Printf("ServiceHandler error:List %v", err)
		return ctx.ServerError()
	}
	return ctx.Success(services)
}

func (s *ServiceHandler) New(ctx Context, f form.NewService) error {
	switch db.Services.Create(ctx.Request().Context(), &db.NewServiceOption{
		Name:    f.Name,
		Health:  f.Health,
		GroupID: f.GroupID,
	}) {
	case nil:
		return ctx.Success()
	case db.ErrServiceGroupNotExists:
		return ctx.Error(400, "Service Group ID not exists.")
	case db.ErrServiceAlreadyExists:
		return ctx.Error(400, "Service in group already exists.")
	default:
		log.Println("Service.New: unknown error.")
		return ctx.ServerError()
	}
}
