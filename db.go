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

// https://github.com/jinzhu/gorm/issues/146
type DBLintStatus struct {
	gorm.Model
	Id           int64
	Organization string
	Repository   string
	Branch       string
	Status       bool
}
