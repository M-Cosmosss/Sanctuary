package route

import (
	"github.com/pkg/errors"
	"net/http"
)

type Table map[string]*Service

type methodTables map[string]Table

type ServicesMap map[string]*Service

type MainRouter struct {
	RTables          methodTables
	WTables          methodTables
	AllowedMethods   []string
	GlobalMiddleware HandlersChain
}

func New() *MainRouter {

}

func (m *MainRouter) AddRoute(method string, path string, service *Service) error {
	var err error
	err = m.RTables.addRoute(method, path, service)
	if err != nil {
		return err
	}
	err = m.WTables.addRoute(method, path, service)
	if err != nil {
		return err
	}
	return nil
}

func (m methodTables) addRoute(method string, path string, service *Service) error {
	table, ok := m[method]
	if !ok {
		m[method] = make(Table)
		table = m[method]
	}
	if _, ok := table[path]; ok {
		return errors.New("redefine route")
	} else {
		table[path] = service
		return nil
	}
}

func (m *MainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if t, ok := m.RTables[method]; ok {
		if s, ok := t[r.URL.String()]; ok {
			h := combineHandlers(s.Group.Middlewares, s.Handlers)
			c := newContext(w, r, h)
			c.Next()
		} else {
			panic("502")
		}
	} else {
		panic("unsupported method")
	}
}

func combineHandlers(chain HandlersChain, chains ...HandlersChain) HandlersChain {
	for _, c := range chains {
		chain = append(chain, c...)
	}
	return chain
}



type Plugin struct {
}

type ServiceGroup struct {
	GroupName   string
	Services    ServicesMap
	Middlewares HandlersChain
}

type Service struct {
	ServiceName string
	ServiceType int
	Handlers    HandlersChain
	Group       *ServiceGroup
}

const (
	ReverseProxy = 1
)

//func (receiver Service) name() {
//	http.ListenAndServe()
//}
