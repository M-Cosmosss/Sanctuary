package gateway

import (
	"log"
	"net/http"
	"sanctuary/internal/db"
	"sanctuary/internal/proxy"
)

func ServiceHandler(ctx *Context) {
	s, err := db.Services.GetServiceFromCache(ctx.ServiceID)
	if err != nil {
		log.Printf("ServiceHandler :get service error.serviceID: %d", ctx.ServiceID)
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(s.Nodes) == 0 {
		log.Printf("ServiceHandler :Don't have any avaliable service node.serviceID: %d", ctx.ServiceID)
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	proxy.ReverseProxy(ctx.ResponseWriter, ctx.Req, s.Nodes[0])
}
