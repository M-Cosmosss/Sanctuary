package gateway

import (
	"context"
	"log"
	"net/http"
	"sanctuary/internal/db"
	"time"
)

var mr *MainRouter

var IsSync bool

type Table map[string]*Route

type methodTables map[string]Table

type Route struct {
	RouteGroupID uint
	ServiceID    uint
}

type task struct {
	Type      int
	Parameter interface{}
	Ctx       context.Context
	ErrorCh   chan error
}

type MainRouter struct {
	Tables           methodTables
	AllowedMethods   []string
	GlobalMiddleware HandlersChain
}

type Gateway interface {
	AddRoute()
	UpdateRoute()
	DeleteRoute()
}

type NewMainRouterOption struct {
	Methods []string
}

func NewMainRouter(o *NewMainRouterOption) *MainRouter {
	m := &MainRouter{
		Tables:           buildTableFromDB(),
		AllowedMethods:   o.Methods,
		GlobalMiddleware: []Handler{ServiceHandler},
	}
	return m
}

func Run() {
	mr = NewMainRouter(&NewMainRouterOption{
		Methods: []string{http.MethodGet, http.MethodPost},
	})
	go mr.syncWorker()
	http.ListenAndServe("127.0.0.1:2832", mr)
}

func (m *MainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if t, ok := m.Tables[method]; ok {
		if route, ok := t[r.URL.String()]; ok {
			c := newContext(w, r, m.GlobalMiddleware, route)
			c.Next()
		} else {
			http.NotFound(w, r)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func buildTableFromDB() methodTables {
	t := make(methodTables)
	routes, err := db.Routes.GetAllOrderByMethod(context.Background())
	if err != nil {
		log.Fatalf("build table: %v", err)
	}
	for _, r := range routes {
		if _, ok := t[r.Method]; !ok {
			t[r.Method] = make(Table)
		}
		t[r.Method][r.FullPath] = &Route{
			RouteGroupID: r.GroupID,
			ServiceID:    r.ServiceID,
		}
	}
	return t
}

func (m *MainRouter) syncWorker() {
	t := time.NewTicker(time.Second * 4)
	for {
		select {
		case <-t.C:
			log.Println("Main Router synced.")
			m.Tables = buildTableFromDB()
		}
	}
}

//type ServiceGroup struct {
//	GroupName   string
//	Services    ServicesMap
//	Middlewares HandlersChain
//}

type Service struct {
	Name   string
	Group  string
	Plugin string
}
