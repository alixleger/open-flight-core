package db

import (
	_ "github.com/jinzhu/gorm" // For annotations
)

// FavFlight model
type FavFlight struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	UserID      uint    `json:"user"`
	FlightID    uint    `json:"flight"`
	TargetPrice float32 `sql:"type:decimal(10,2);" json:"target_price"`
}
