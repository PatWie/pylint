package pylint

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DBInstallation struct {
	gorm.Model
	Installation int64
	Sender       int64
}

func CreateInstallation(cc InstallationPayload) {
	Database.Create(&DBInstallation{Sender: cc.Sender.ID,
		Installation: cc.Installation.ID})
}

func DeleteInstallation(cc InstallationPayload) {
	Database.Where("installation = ?", cc.Installation.ID).Delete(DBInstallation{})
}

// https://github.com/jinzhu/gorm/issues/146
type DBLintStatus struct {
	gorm.Model
	Organization string
	Repository   string
	Branch       string
	Commit       string
	Msg          string
	Status       int
}

// wrapper for database
var Database *gorm.DB

func ConnectDatabase(cfg Config) error {
	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	var err error
	Database, err = gorm.Open("postgres", connInfo)

	return err
}
