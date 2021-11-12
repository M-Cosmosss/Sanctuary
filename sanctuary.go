package main

import (
	"sanctuary/internal/route"
	"sanctuary/internal/service"
)

func main(){
	//gin.Default()
	//r:=route.New()
	//r.Use(logg.Log)
	//e99:=r.Group("/e99")
	//e99.GET("/a", func(context *route.Context) {
	//	proxy.ReverseProxy("http://0:9091",context)
	//})
	//err=r.GET("/e99/a")
	//http.Handle()
	//http.ListenAndServe(":6788",r)
	//service.InitServiceCenter()
	center:=service.NewServiceCenterBasedOnTextFile(&service.TextFilePersistenceOption{FilePath: "./services.json"})
	s:=&service.Service{
		ServiceName: "echoo",
		ServiceType: route.ReverseProxy,
		Handlers:    nil,
		GroupName:   "e999",
	}
	sg:=&service.ServiceGroup{
		GroupName:   "e999",
		Services:    make(service.ServicesMap),
		Middlewares: nil,
	}
	sg.Services[s.ServiceName]=s
	center.NewGroup(sg)
	center.Store()
	//service.Build("./service.json")
}




//func (r *route.Router){
//
//}