package skyscanner

import "time"

// Flight ressource
type Flight struct {
	Carrier          string
	Price            float64
	OriginPlace      Place
	DestinationPlace Place
	DepartureDate    time.Time
}
