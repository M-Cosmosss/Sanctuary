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
			f.Get("/route", form.Bind(form.GetRouteOption{}), route.List)
			f.Post("/route", form.Bind(form.NewRoute{}), route.New)
			f.Put("/route")
			f.Delete("/route")
			f.Get("/group", routeGroup.List)
			f.Post("/group", form.Bind(form.NewRouteGroup{}), routeGroup.New)
			f.Put("/group")
			f.Delete("/group")
		})

		f.Group("/center", func() {
			f.Get("/service", form.Bind(form.GetServiceOption{}), service.List)
			f.Post("/group")
			f.Post("/service", form.Bind(form.NewService{}), service.New)
			f.Post("/register", form.Bind(form.NewService{}), service.New)
		})
	})
	f.Run()
}

func Run() {

}
