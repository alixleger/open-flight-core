package db

import (
	"time"

	_ "github.com/jinzhu/gorm" // For annotations
)

// Flight model
type Flight struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	DepartPlaceID  uint      `json:"depart_place_id"`
	DepartPlace    Place     `gorm:"foreignkey:DepartPlaceID" json:"depart_place"`
	ArrivalPlaceID uint      `json:"arrival_place_id"`
	ArrivalPlace   Place     `gorm:"foreignkey:ArrivalPlaceID" json:"arrival_place"`
	DepartDate     time.Time `json:"depart_date"`
	ArrivalDate    time.Time `json:"arrival_date"`
	ExternalID     string    `json:"external_id"`
	Price          uint      `json:"price"`
}
