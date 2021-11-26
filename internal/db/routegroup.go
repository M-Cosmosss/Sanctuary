package db

import (
	"context"
	"log"
	"sanctuary/internal/utils"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RouteGroup struct {
	gorm.Model
	Name   string         `gorm:"unique"`
	Path   string         `gorm:"unique"`
	Plugin pq.StringArray `gorm:"type:varchar(255)[]"`
}

var RouteGroups RouteGroupsStore
var _ RouteGroupsStore = &routeGroups{}

type RouteGroupsStore interface {
	Create(ctx context.Context, g *NewRouteGroupOption) error
	//UpdateRouteGroup()
	//DeleteRouteGroup()

	Get(ctx context.Context) ([]*RouteGroup, error)
	GetByID(ctx context.Context, id uint) (*RouteGroup, error)
	GetByName(ctx context.Context, n string) (*RouteGroup, error)
}

type routeGroups struct {
	*gorm.DB
}

func NewRouteGroupsStore(db *gorm.DB) *routeGroups {
	return &routeGroups{DB: db}
}

type NewRouteGroupOption RouteGroup

func (db *routeGroups) Create(ctx context.Context, g *NewRouteGroupOption) error {
	if err := db.WithContext(ctx).Create((*RouteGroup)(g)).Error; err != nil {
		if utils.IsUniqueError(err, "route_groups_name_key") {
			log.Println("Create: repeat route groups name")
			return ErrRouteGroupAlreadyExists
		}
		return ErrUnknown
	}
	return nil
}

func (db *routeGroups) UpdateRouteGroup() {
	panic("implement me")
}

func (db *routeGroups) DeleteRouteGroup() {
	panic("implement me")
}

func (db *routeGroups) Get(ctx context.Context) ([]*RouteGroup, error) {
	var routeGroups []*RouteGroup
	return routeGroups, db.WithContext(ctx).Find(&routeGroups).Error
}

func (db *routeGroups) GetByID(ctx context.Context, id uint) (*RouteGroup, error) {
	var g RouteGroup
	switch db.WithContext(ctx).Model(&RouteGroup{}).Where("id = ?", id).First(&g).Error {
	case nil:
		return &g, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrRouteGroupNotExists
	default:
		return nil, ErrUnknown
	}
}

func (db *routeGroups) GetByName(ctx context.Context, n string) (*RouteGroup, error) {
	var g RouteGroup
	switch db.WithContext(ctx).Model(&RouteGroup{}).Where("name = ?", n).First(&g).Error {
	case nil:
		return &g, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrRouteGroupNotExists
	default:
		return nil, ErrUnknown
	}
}
