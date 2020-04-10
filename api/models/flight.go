package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Flight model
type Flight struct {
	gorm.Model
	CompanyID        uint      `json:"company"`
	DepartAirportID  uint      `json:"depart_airport"`
	ArrivalAirportID uint      `json:"arrival_airport"`
	DepartDate       time.Time `json:"depart_date"`
	ArrivalDate      time.Time `json:"arrival_date"`
}
