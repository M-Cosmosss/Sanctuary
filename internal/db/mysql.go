package db

import (
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Mysql *gorm.DB

func MysqlInit() error {
	db, err := gorm.Open(mysql.Open("root:mysql123@tcp(127.0.0.1:3306)/sanctuary?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		return errors.Wrap(err, "mysql init")
	}
	err = db.AutoMigrate(&RequestLog{})
	if err != nil {
		return errors.Wrap(err, "mysql autoMigrate")
	}

	sql, _ := db.DB()
	sql.SetMaxOpenConns(300)
	Mysql = db
	return nil
}

//func
