package backend

import (
	"sanctuary/internal/route"
	"sanctuary/internal/servicecenter"

	"gopkg.in/macaron.v1"
)

var r route.MainRouter
var c servicecenter.ServiceCenter

func Run() {
	m := macaron.Classic()

}


func