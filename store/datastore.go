// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/patwie/pylint/model"
	"log"
)

const (
	DriverSqlite = "sqlite3"
)

type Datastore interface {
	CreateInstallation(*model.Installation) error
	GetInstallation(installationID int) (*model.Installation, error)
	DeleteInstallation(installationID int) error
}

type datastore struct {
	db *gorm.DB
}

var ds = newDatastore()

func DS() Datastore { return ds }

func newDatastore() Datastore {

	config := model.GetConfiguration()

	db, err := gorm.Open(DriverSqlite, config.Pylint.DatabaseFile)
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}

	db.AutoMigrate(&model.Installation{})

	return &datastore{db: db}
}
