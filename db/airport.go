package db

import (
	"github.com/jinzhu/gorm"
)

// Airport model
type Airport struct {
	gorm.Model
	Name string `json:"name"`
	City string `json:"city"`
}
