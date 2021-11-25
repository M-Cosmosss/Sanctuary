package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var Services ServicesStore
var _ ServicesStore = &services{}

func NewServicesStore(db *gorm.DB) *services {
	return &services{DB: db}
}

type Service struct {
	gorm.Model
	Name    string
	Nodes   pq.StringArray `gorm:"type:varchar(255)[]"`
	Health  pq.StringArray `gorm:"type:varchar(255)[]"`
	GroupID uint
}

type ServicesStore interface {
}
type services struct {
	*gorm.DB
}

func (db *services) NewService() {

}
