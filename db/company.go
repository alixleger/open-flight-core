package db

import (
	_ "github.com/jinzhu/gorm" // For annotations
)

// Company model
type Company struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `gorm:"size:32;not null;unique_index" json:"name"`
}
