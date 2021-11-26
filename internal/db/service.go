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
	AddNode(ctx context.Context, id uint, url string) error
	Create(ctx context.Context, o *NewServiceOption) error
	Get(ctx context.Context, o *GetServiceOption) ([]*Service, error)
	GetAll(ctx context.Context) ([]*Service, error)
	GetByID(ctx context.Context, id uint) (*Service, error)
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

func (db *services) AddNode(ctx context.Context, id uint, url string) error {
	s, err := db.GetByID(ctx, id)
	if err != nil {
		return err
	}
	s.Nodes = append(s.Nodes, url)
	switch db.WithContext(ctx).Where("id = ?", id).Update("nodes", s.Nodes).Error {
	case nil:
		return nil
	default:
		log.Fatalf("AddNode error:%v", err)
		return err
	}
}

type GetServiceOption struct {
	OrderBy  string
	Page     int
	PageSize int
}

func (db *services) Get(ctx context.Context, o *GetServiceOption) ([]*Service, error) {
	var services []*Service
	if o.OrderBy == "" {
		o.OrderBy = "group_id ASC,name ASC"
	}
	if o.Page <= 0 {
		o.Page = 1
	}
	if o.PageSize <= 0 {
		o.PageSize = 50
	}
	return services, db.WithContext(ctx).
		Offset((o.Page - 1) * o.PageSize).Limit(o.PageSize).Order(o.OrderBy).
		Find(&services).Error
}

func (db *services) GetAll(ctx context.Context) ([]*Service, error) {
	var s []*Service
	return s, db.WithContext(ctx).Find(&s).Error
}

func (db *services) GetByID(ctx context.Context, id uint) (*Service, error) {
	var s *Service
	switch db.WithContext(ctx).Where("id = ?", id).First(s).Error {
	case nil:
		return s, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrServiceNotExists
	default:
		return nil, ErrUnknown
	}
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
