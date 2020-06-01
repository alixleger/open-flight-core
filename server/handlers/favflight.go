package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type favFlightsType struct {
	FavFlights []models.FavFlight `json:"fav_flights"`
}

type favFlightInput struct {
	FlightID    uint `json:"flight_id" validate:"required"`
	TargetPrice uint `json:"target_price" validate:"required"`
}

// GetFavFlights function
func GetFavFlights(c *gin.Context) {
	userInterface, _ := c.Get(IdentityKey)
	user := userInterface.(*models.User)
	var res favFlightsType
	var favFlights []models.FavFlight
	db := c.MustGet("db").(*gorm.DB)

	db.Where("user_id = ?", user.ID).Find(&favFlights).Preload("Flight")
	res.FavFlights = favFlights

	c.JSON(http.StatusOK, res)
	return
}

// PostFavFlight function
func PostFavFlight(c *gin.Context) {
	userInterface, _ := c.Get(IdentityKey)
	user := userInterface.(*models.User)

	// Parse and check post data
	var input favFlightInput
	validator := validator.New()
	if err := c.ShouldBindJSON(&input); err != nil || validator.Struct(input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should have correct flight_id and target_price fields"})
		return
	}

	var flight models.Flight
	db := c.MustGet("db").(*gorm.DB)
	db.First(&flight, input.FlightID)

	if input.FlightID != flight.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This flight does not exist"})
		return
	}

	if flight.Price < input.TargetPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current flight price is already lower than that target price"})
		return
	}

	// Create favFlight entity
	var favFlight models.FavFlight
	favFlight.FlightID = flight.ID
	favFlight.UserID = user.ID
	favFlight.TargetPrice = input.TargetPrice
	db.Create(&favFlight)

	// Insert flight price in influxdb
	influxClient := *c.MustGet("influxdbClient").(*client.Client)
	err := InsertFlightPrice(influxClient, user.Email, flight.ExternalID, flight.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// InsertFlightPrice function
func InsertFlightPrice(influxClient client.Client, userID string, flightID string, price uint) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  os.Getenv("INFLUXDB_DATABASE"),
		Precision: "h",
	})
	if err != nil {
		log.Println(err)
		return err
	}

	tags := map[string]string{"user_id": userID, "flight": flightID}
	fields := map[string]interface{}{"price": price}

	pt, err := client.NewPoint("flightprices", tags, fields, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	bp.AddPoint(pt)

	if err := influxClient.Write(bp); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
