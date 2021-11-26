package backend

import (
	"sanctuary/internal/form"

	"github.com/flamego/flamego"
)

var route = NewRouteHandler()
var routeGroup = NewRouteGroupHandler()
var service = NewServiceHandler()
var serviceGroup = NewServiceGroupHandler()

func NewRouter() {
	f := flamego.Classic()
	f.Use(Contexter())
	f.Group("/api", func() {

		f.Group("/gateway", func() {
			f.Combo("/route").
				Get(form.Bind(form.GetRouteOption{}), route.List).
				Post(form.Bind(form.NewRoute{}), route.New).
				Put().
				Delete()
			f.Combo("/group").
				Get(routeGroup.List).
				Post(form.Bind(form.NewRouteGroup{}), routeGroup.New).
				Put().
				Delete()
		})

		f.Group("/center", func() {
			f.Combo("/service").
				Get(form.Bind(form.GetServiceOption{}), service.List).
				Post(form.Bind(form.NewService{}), service.New)
			f.Post("/register", form.Bind(form.NewServiceNode{}), service.RegisterNode)
			f.Post("/group", form.Bind(form.NewServiceGroup{}), serviceGroup.New)
		})
	})
	go f.Run()
}

func Run() {

}
