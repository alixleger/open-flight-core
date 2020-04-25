package db

import (
	"time"

	_ "github.com/jinzhu/gorm" // For annotations
)

// Flight model
type Flight struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	CompanyID      uint      `json:"company_id"`
	DepartAirport  string    `json:"depart_airport"`
	ArrivalAirport string    `json:"arrival_airport"`
	DepartDate     time.Time `json:"depart_date"`
	ArrivalDate    time.Time `json:"arrival_date"`
	ExternalID     string    `json:"external_id"`
	Price          uint      `json:"price"`
}
