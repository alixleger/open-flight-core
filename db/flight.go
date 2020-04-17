package db

import (
	"time"

	_ "github.com/jinzhu/gorm" // For annotations
)

// Flight model
type Flight struct {
	ID               uint      `gorm:"primary_key" json:"id"`
	CompanyID        uint      `json:"company"`
	DepartAirportID  uint      `json:"depart_airport"`
	ArrivalAirportID uint      `json:"arrival_airport"`
	DepartDate       time.Time `json:"depart_date"`
	ArrivalDate      time.Time `json:"arrival_date"`
}
