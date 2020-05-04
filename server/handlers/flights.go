package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/alixleger/open-flight-core/services/skyscanner"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type flightsType struct {
	Flights []models.Flight `json:"flights"`
}

// GetFlights function
func GetFlights(c *gin.Context) {
	departPlaceID, exists := c.GetQuery("departPlaceID")
	departPlaceIDValue, err := strconv.Atoi(departPlaceID)
	if !exists || departPlaceID == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flights endpoint must have an integer departPlaceID param"})
		return
	}

	arrivalPlaceID, exists := c.GetQuery("arrivalPlaceID")
	arrivalPlaceIDValue, err := strconv.Atoi(arrivalPlaceID)
	if !exists || arrivalPlaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flights endpoint must have an integer arrivalPlaceID param"})
		return
	}

	departDate, exists := c.GetQuery("departDate")
	departDateValue, err := strconv.ParseInt(departDate, 10, 64)
	if !exists || departDate == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flights endpoint must have an integer departDate param (timestamp)"})
		return
	}
	departDateTimeValue := time.Unix(departDateValue, 0)

	var flights []models.Flight
	db := c.MustGet("db").(*gorm.DB)
	db.Where("depart_place_id = ? AND arrival_place_id = ? AND depart_date = ?", departPlaceIDValue, arrivalPlaceIDValue, departDateTimeValue).Find(&flights)

	if len(flights) == 0 {
		var departPlace models.Place
		var arrivalPlace models.Place

		db.First(&departPlace, departPlaceIDValue)
		db.First(&arrivalPlace, arrivalPlaceIDValue)

		if departPlace.ExternalID == "" || arrivalPlace.ExternalID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "departDate and arrivalDate must exist"})
			return
		}

		skyscannerClient := c.MustGet("skyscannerClient").(*skyscanner.Client)
		skyscannerFlights := skyscannerClient.GetFlights(departPlace.ExternalID, arrivalPlace.ExternalID, fmt.Sprintf("%d-%02d-%02d", departDateTimeValue.Year(), departDateTimeValue.Month(), departDateTimeValue.Day()))
		for _, flight := range skyscannerFlights {
			var flightDepartPlace models.Place
			var flightArrivalPlace models.Place

			db.Where("external_id = ?", flight.OriginPlace.PlaceId).First(&flightDepartPlace)
			db.Where("external_id = ?", flight.DestinationPlace.PlaceId).First(&flightArrivalPlace)

			if flightDepartPlace.ID <= 0 {
				db.Create(&flightDepartPlace)
			}
			if flightArrivalPlace.ID <= 0 {
				db.Create(&flightArrivalPlace)
			}

			dbFlight := models.Flight{
				DepartDate:     flight.DepartureDate,
				Price:          uint(flight.Price),
				ExternalID:     fmt.Sprintf("%s-%s-%s-%d", flight.Carrier, flight.OriginPlace.PlaceId, flight.DestinationPlace.PlaceId, flight.DepartureDate.Unix()),
				DepartPlaceID:  flightDepartPlace.ID,
				DepartPlace:    flightDepartPlace,
				ArrivalPlaceID: flightArrivalPlace.ID,
				ArrivalPlace:   flightArrivalPlace,
			}
			db.Create(&dbFlight)
			flights = append(flights, dbFlight)
		}
	}

	var res flightsType
	res.Flights = flights

	c.JSON(http.StatusOK, res)
	return
}
