package main

import (
	"sanctuary/internal/backend"
	"sanctuary/internal/db"
	"sanctuary/internal/gateway"
)

func main() {
	//var err error
	db.Init()
	backend.NewRouter()
	gateway.Run()
	//err = db.RouteGroups.Create(&db.NewRouteGroupOption{
	//	Name:   "group1",
	//	Path:   "/1",
	//	Plugin: nil,
	//})
	//if err != nil {
	//	log.Fatalf("1:%v", err)
	//}
	//err = db.RouteGroups.Create(&db.NewRouteGroupOption{
	//	Name:   "group1",
	//	Path:   "/2",
	//	Plugin: nil,
	//})
	//if err != nil {
	//	log.Fatalf("2:%v", err)
	//}
}
