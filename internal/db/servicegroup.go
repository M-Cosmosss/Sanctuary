package db

import (
	"context"
	"log"
	"sanctuary/internal/utils"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var ServiceGroups ServiceGroupStore
var _ ServiceGroupStore = &serviceGroups{}

type ServiceGroupStore interface {
	Create(ctx context.Context, g *NewServiceGroupOption) error
	//UpdateServiceGroup()
	//DeleteServiceGroup()
	//GetByName(n string) (*ServiceGroup, error)
	GetByID(ctx context.Context, id uint) (*ServiceGroup, error)
}

type serviceGroups struct {
	*gorm.DB
}

func NewServiceGroupsStore(db *gorm.DB) *serviceGroups {
	return &serviceGroups{DB: db}
}

type ServiceGroup struct {
	gorm.Model
	Name   string         `gorm:"unique"`
	Plugin pq.StringArray `gorm:"type:varchar(255)[]"`
}

type NewServiceGroupOption ServiceGroup

func (db *serviceGroups) Create(ctx context.Context, g *NewServiceGroupOption) error {
	if err := db.WithContext(ctx).Create((*ServiceGroup)(g)).Error; err != nil {
		if utils.IsUniqueError(err, "service_groups_name_key") {
			return ErrServiceGroupAlreadyExists
		}
		return ErrUnknown
	}
	return nil
}

func (db *serviceGroups) GetByID(ctx context.Context, id uint) (*ServiceGroup, error) {
	var sg *ServiceGroup
	switch db.WithContext(ctx).Where("id = ?", id).First(sg).Error {
	case nil:
		return sg, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrServiceGroupNotExists
	default:
		log.Println("ServiceGroups.GetByID: unknown error.")
		return nil, ErrUnknown
	}
}
