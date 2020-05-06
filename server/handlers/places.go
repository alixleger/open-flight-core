package handlers

import (
	"net/http"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/alixleger/open-flight-core/services/skyscanner"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type placesType struct {
	Places []models.Place `json:"places"`
}

// GetPlaces function
func GetPlaces(c *gin.Context) {
	queryString, exists := c.GetQuery("query")
	if !exists || queryString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "places endpoint must have a query param"})
		return
	}
	var places []models.Place
	db := c.MustGet("db").(*gorm.DB)
	db.Where("name LIKE ?", "%"+queryString+"%").Or("country LIKE ?", "%"+queryString+"%").Find(&places)

	if len(places) == 0 {
		skyscannerClient := c.MustGet("skyscannerClient").(*skyscanner.Client)
		skyscannerPlaces, err := skyscannerClient.GetPlaces(queryString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "External services error"})
			return
		}

		for _, place := range skyscannerPlaces {
			dbPlace := models.Place{
				ExternalID: place.PlaceId,
				Name:       place.PlaceName,
				Country:    place.CountryName,
			}
			db.Create(&dbPlace)
			places = append(places, dbPlace)
		}
	}

	var res placesType
	res.Places = places

	c.JSON(http.StatusOK, res)
	return
}
