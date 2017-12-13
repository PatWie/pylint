package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBInstallation struct {
	gorm.Model
	Installation int64
	Sender       int64
}
