package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var AllTables = []interface{}{
	Route{},
	RouteGroup{},
	Service{},
	ServiceGroup{},
}

// Init initializes the database.
func Init() error {
	dsn := "user=postgres password=pgpg123 dbname=sanctuary port=5432 sslmode=disable host=127.0.0.1 TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "open connection")
	}

	// Migrate databases.
	if db.AutoMigrate(AllTables...) != nil {
		return errors.Wrap(err, "auto migrate")
	}

	SetDatabaseStore(db)
	TableInit()
	return nil
}

func SetDatabaseStore(db *gorm.DB) {
	Routes = NewRouteStore(db)
	RouteGroups = NewRouteGroupsStore(db)
	Services = NewServicesStore(db)
	ServiceGroups = NewServiceGroupsStore(db)
}

func TableInit() {
	RouteGroups.Create(context.Background(), &NewRouteGroupOption{
		Name:   "Default",
		Path:   "/",
		Plugin: nil,
	})
	ServiceGroups.Create(context.Background(), &NewServiceGroupOption{
		Name:   "Default",
		Plugin: nil,
	})
}
