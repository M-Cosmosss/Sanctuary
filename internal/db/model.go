package db

import "gorm.io/gorm"

type RequestLog struct {
	gorm.Model
	IP      string
	Method  string
	Path    string
	Success bool
	ErrMsg  string
}
