// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package model

import (
	"github.com/jinzhu/gorm"
)

type Installation struct {
	gorm.Model
	SenderID       int
	InstallationID int
}
