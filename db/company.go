package db

import (
	"github.com/jinzhu/gorm"
)

// Company model
type Company struct {
	gorm.Model
	Name string `gorm:"size:32;not null;unique_index" json:"name"`
}
