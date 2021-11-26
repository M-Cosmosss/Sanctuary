package db

import (
	"context"
	"log"
	"time"

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
	Create(ctx context.Context, o *NewServiceOption) error
	GetAll(ctx context.Context) ([]*Service, error)
	GetByNameAndGroupID(ctx context.Context, id uint, n string) (*Service, error)
	GetServiceFromCache(id uint) (*Service, error)
}
type services struct {
	*gorm.DB
	cache  map[uint]*Service
	IsSync bool
}

type NewServiceOption struct {
	Name    string
	Health  []string
	GroupID uint
}

func (db *services) Create(ctx context.Context, o *NewServiceOption) error {
	s := &Service{
		Name:    o.Name,
		Nodes:   nil,
		Health:  o.Health,
		GroupID: o.GroupID,
	}
	if _, err := ServiceGroups.GetByID(ctx, o.GroupID); err != nil {
		return err
	}
	if _, err := Services.GetByNameAndGroupID(ctx, o.GroupID, o.Name); err != nil {
		if err != ErrServiceNotExists {
			return err
		}
	} else {
		return ErrServiceAlreadyExists
	}
	switch db.WithContext(ctx).Create(s).Error {
	case nil:
		db.IsSync = false
		return nil
	default:
		log.Println("Services.Create: unknown error.")
		return ErrUnknown
	}
}

func (db *services) GetAll(ctx context.Context) ([]*Service, error) {
	var s []*Service
	return s, db.WithContext(ctx).Find(&s).Error
}

func (db *services) GetByNameAndGroupID(ctx context.Context, id uint, n string) (*Service, error) {
	var s *Service
	switch db.WithContext(ctx).Where("group_id = ? AND name = ?", id, n).Error {
	case nil:
		return s, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrServiceNotExists
	default:
		return nil, ErrUnknown
	}
}

func (db *services) GetServiceFromCache(id uint) (*Service, error) {
	s, ok := db.cache[id]
	if !ok {
		return nil, ErrServiceNotExists
	}
	return s, nil
}

func (db *services) Syncer() {
	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-t.C:
			if !db.IsSync {
				db.syncCache()
				db.IsSync = true
			}
		}
	}
}

func (db *services) syncCache() {
	s, err := Services.GetAll(context.Background())
	if err != nil {
		log.Fatalf("Services.syncCache: GetAll error.")
	}
	m := map[uint]*Service{}
	for _, service := range s {
		m[service.ID] = service
	}
	db.cache = m
	log.Println("Services cache synced")
}
