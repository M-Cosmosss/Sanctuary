package db

import (
	"context"
	"path"
	"sanctuary/internal/utils"

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
	GetAllOrderByMethod(ctx context.Context) ([]*Route, error)
	GetByGroupID(ctx context.Context, gid uint) ([]*Route, error)
	GetByMethodAndPath(ctx context.Context, method string, fullPath string) (*Route, error)
	//DeleteRoute(DeleteRouteOption) error
	//UpdateRoute(UpdateRouteOption) error
}

type NewRouteOption struct {
	Method    string
	Path      string
	GroupID   uint
	ServiceID uint
}

func (db *routes) Create(ctx context.Context, o *NewRouteOption) error {
	if !utils.IsHTTPMethod(o.Method) {
		return ErrNotHTTPMethod
	}
	g, err := RouteGroups.GetByID(ctx, o.GroupID)
	if err != nil {
		return err
	}
	r := &Route{
		Method:    o.Method,
		FullPath:  path.Join(g.Path, o.Path),
		Path:      o.Path,
		GroupID:   o.GroupID,
		ServiceID: o.ServiceID,
	}
	if _, err = db.GetByMethodAndPath(ctx, r.Method, r.FullPath); err != ErrRouteNotExists {
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

func (db *routes) GetAllOrderByMethod(ctx context.Context) ([]*Route, error) {
	var routes []*Route
	return routes, db.WithContext(ctx).Order("method").Find(&routes).Error
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
	var r Route
	switch db.WithContext(ctx).Model(&Route{}).Where("full_path = ? AND method = ?", fullPath, method).First(&r).Error {
	case nil:
		return &r, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrRouteNotExists
	default:
		return nil, ErrUnknown
	}
}
