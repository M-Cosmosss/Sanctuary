package backend

import (
	"sanctuary/internal/form"
	"sanctuary/internal/sc"
)

type CenterHandler struct{}

func NewCenterHandler() *CenterHandler {
	return &CenterHandler{}
}

func (*CenterHandler) New() {

}

func (*CenterHandler) ListAll(ctx Context) error {
	m, _ := sc.Sc.Get()
	return ctx.Success(m)
}

func (*CenterHandler) NewGroup(ctx Context, f form.NewServiceGroup) error {
	err := sc.Sc.RegisterServiceGroup(&sc.ServiceGroup{Name: f.Name})
	if err != nil {
		return ctx.Error(40100, "group already exists")
	}
	return ctx.Success()
}

//func (*CenterHandler) NewService(ctx Context, f form.NewService) error {
//	err := sc.Sc.RegisterService(&sc.Service{
//		Name:      f.Name,
//		GroupName: f.Group,
//	})
//	if err != nil {
//		return ctx.Error(40100, "service already exists")
//	}
//	return ctx.Success()
//}
