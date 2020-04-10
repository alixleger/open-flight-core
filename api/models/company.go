package models

import (
	"github.com/jinzhu/gorm"
)

// Company model
type Company struct {
	gorm.Model
	Name string `json:"name"`
}
