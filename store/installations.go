// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package store

import (
	"github.com/patwie/pylint/model"
)

func (ds *datastore) CreateInstallation(i *model.Installation) error {
	return ds.db.Create(i).Error
}

func (ds *datastore) GetInstallation(installationID int) (*model.Installation, error) {
	data := new(model.Installation)
	err := ds.db.Where("installation_id = ?", installationID).First(&data).Error
	return data, err

}

func (ds *datastore) DeleteInstallation(installationID int) error {
	return ds.db.Where("installation_id = ?", installationID).Delete(model.Installation{}).Error

}
