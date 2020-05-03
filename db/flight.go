package db

import (
	"time"

	_ "github.com/jinzhu/gorm" // For annotations
)

// Flight model
type Flight struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	DepartPlaceID  uint      `json:"depart_place_id"`
	ArrivalPlaceID uint      `json:"arrival_place_id"`
	DepartDate     time.Time `json:"depart_date"`
	ArrivalDate    time.Time `json:"arrival_date"`
	ExternalID     string    `json:"external_id"`
	Price          uint      `json:"price"`
}
