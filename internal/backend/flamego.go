package backend

import (
	"sanctuary/internal/form"

	"github.com/flamego/flamego"
)

var center = NewCenterHandler()
var route = NewRouteHandler()
var routeGroup = NewRouteGroupHandler()

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
			f.Get("/service", center.ListAll)
			f.Post("/group")
			f.Post("/service")
		})
	})
	f.Run()
}

func Run() {

}
