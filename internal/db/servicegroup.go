package db

import (
	"sanctuary/internal/utils"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var ServiceGroups ServiceGroupStore
var _ ServiceGroupStore = &serviceGroups{}

type ServiceGroupStore interface {
	NewServiceGroup(g *NewServiceGroupOption) error
	//UpdateServiceGroup()
	//DeleteServiceGroup()
	//GetByName(n string) (*ServiceGroup, error)
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

func (db *serviceGroups) NewServiceGroup(g *NewServiceGroupOption) error {
	if err := db.Create((*ServiceGroup)(g)).Error; err != nil {
		if utils.IsUniqueError(err, "service_groups_name_key") {
			return ErrServiceGroupAlreadyExists
		}
		return ErrUnknown
	}
	return nil
}
