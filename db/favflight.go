package db

import (
	_ "github.com/jinzhu/gorm" // For annotations
)

// FavFlight model
type FavFlight struct {
	ID          uint `gorm:"primary_key" json:"id"`
	UserID      uint `json:"user_id"`
	FlightID    uint `json:"flight_id"`
	TargetPrice uint `json:"target_price"`
}
