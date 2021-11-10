package main

import (
	"net/http"
	logg "sanctuary/internal/log"
	"sanctuary/internal/proxy"
	"sanctuary/internal/route"
)

func main(){
	var err error
	r:=route.New()
	r.Use(logg.Log)
	e99:=r.Group("/e99")
	e99.GET("/a", func(context *route.Context) {
		proxy.ReverseProxy("http://0:9091",context)
	})
	err=r.GET("/e99/a")

	http.ListenAndServe(":6788",r)
}


//func (r *route.Router){
//
//}