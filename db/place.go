package db

import (
	_ "github.com/jinzhu/gorm" // For annotations
)

// Place model
type Place struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	ExternalID string `gorm:"size:32;not null;unique_index" json:"external_id"`
	Name       string `gorm:"size:32;not null;unique_index" json:"name"`
	Country    string `gorm:"size:32;not null" json:"country"`
}
