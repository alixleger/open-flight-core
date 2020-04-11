package db

import (
	"github.com/jinzhu/gorm"
)

// FavFlight model
type FavFlight struct {
	gorm.Model
	UserID   uint `json:"user"`
	FlightID uint `json:"flight"`
}
