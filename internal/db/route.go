package db

import (
	"context"

	"gorm.io/gorm"
)

type Route struct {
	gorm.Model
	Method    string
	FullPath  string `gorm:"index"`
	Path      string
	GroupID   uint `gorm:"index"`
	ServiceID uint `gorm:"index"`
}

var Routes RouteStore
var _ RouteStore = &routes{}

type routes struct {
	*gorm.DB
}

func NewRouteStore(db *gorm.DB) *routes {
	return &routes{DB: db}
}

type RouteStore interface {
	//Load(m *gateway.MainRouter) error
	//Store(m *gateway.MainRouter) error
	Create(ctx context.Context, g *NewRouteOption) error
	Get(ctx context.Context, o *GetRouteOption) ([]*Route, error)
	GetByGroupID(ctx context.Context, gid uint) ([]*Route, error)
	GetByMethodAndPath(ctx context.Context, method string, fullPath string) (*Route, error)
	//DeleteRoute(DeleteRouteOption) error
	//UpdateRoute(UpdateRouteOption) error
}

type NewRouteOption struct {
	Method    string
	FullPath  string
	Path      string
	GroupID   uint
	ServiceID uint
}

func (db *routes) Create(ctx context.Context, o *NewRouteOption) error {
	r := &Route{
		Method:    o.Method,
		FullPath:  o.FullPath,
		Path:      o.Path,
		GroupID:   o.GroupID,
		ServiceID: o.ServiceID,
	}
	if _, err := db.GetByMethodAndPath(ctx, o.Method, o.FullPath); err != ErrRouteNotExists {
		return ErrRouteAlreadyExists
	}

	switch db.WithContext(ctx).Create(r).Error {
	case nil:
		return nil
	default:
		return ErrUnknown
	}
}

type GetRouteOption struct {
	OrderBy  string
	Page     int
	PageSize int
}

func (db *routes) Get(ctx context.Context, o *GetRouteOption) ([]*Route, error) {
	var routes []*Route
	if o.OrderBy == "" {
		o.OrderBy = "group_id ASC,path ASC,method ASC,service_id ASC"
	}
	if o.Page <= 0 {
		o.Page = 1
	}
	if o.PageSize <= 0 {
		o.PageSize = 50
	}
	return routes, db.WithContext(ctx).
		Offset((o.Page - 1) * o.PageSize).Limit(o.PageSize).Order(o.OrderBy).
		Find(&routes).Error
}

func (db *routes) GetByGroupID(ctx context.Context, gid uint) ([]*Route, error) {
	var routes []*Route
	return routes, db.WithContext(ctx).Order("path ASC").Find(&routes).Error
}

func (db *routes) GetByMethodAndPath(ctx context.Context, method string, fullPath string) (*Route, error) {
	r := &Route{}
	switch db.WithContext(ctx).Model(&Route{}).Where("full_path = ? AND method = ?", fullPath, method).First(r).Error {
	case nil:
		return r, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrRouteNotExists
	default:
		return nil, ErrUnknown
	}
}

//func (db *routes) Store(m *RouteGroup) error {
//	if err := db.Where("1=1").Delete(&Route{}).Error; err != nil {
//		panic(err)
//	}
//	routes := make([]*Route, 0)
//	for _, v := range m {
//		for _, r := range v.Routes {
//			routes = append(routes, &Route{
//				Method: r.Method,
//				Path:   r.Path,
//				//Group:        r.Group,
//				//Plugin:       r.Plugin,
//				//Service:      r.Service,
//				//ServiceGroup: r.ServiceGroup,
//			})
//		}
//	}
//	if err := db.Create(routes).Error; err != nil {
//		log.Println("Store error")
//		return err
//	}
//	return nil
//}
//
//func (db *routes) Load(m *gateway.MainRouter) error {
//	//if len(m.RTables) != 0 {
//	//	return errors.New("Must load in a new router")
//	//}
//	//routes := []*Route{}
//	//if err := db.Model(&Route{}).Find(routes).Error; err != nil {
//	//	return err
//	//}
//	//for _, v := range routes {
//	//	if _, ok := m.RouteGroups[v.Group]; !ok {
//	//		m.RouteGroups[v.Group] = &gateway.RouteGroup{Plugin: v.Plugin}
//	//	}
//	//	m.RouteGroups[v.Group].Routes = append(m.RouteGroups[v.Group].Routes, &gateway.Route{
//	//		Method:       v.Method,
//	//		Path:         v.Path,
//	//		Group:        v.Group,
//	//		Plugin:       v.Plugin
//	//		Service:      v.Service,
//	//		ServiceGroup: v.ServiceGroup,
//	//	})
//	//}
//	return nil
//}
