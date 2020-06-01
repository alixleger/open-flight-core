package db

import (
	_ "github.com/jinzhu/gorm" // For annotations
)

// FavFlight model
type FavFlight struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignkey:UserID"`
	FlightID    uint   `json:"flight_id"`
	Flight      Flight `gorm:"foreignkey:FlightID"`
	TargetPrice uint   `json:"target_price"`
}
