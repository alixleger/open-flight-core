package handlers

import (
	"fmt"
	"net/http"
	"os"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
)

// GetFlightPrices function
func GetFlightPrices(c *gin.Context) {
	influxClient := *c.MustGet("influxdbClient").(*client.Client)
	userInterface, _ := c.Get(IdentityKey)
	user := userInterface.(*models.User)

	// Query data
	q := client.Query{
		Command:  fmt.Sprintf("SELECT * FROM flightprices WHERE user_id = '%s'", user.Email),
		Database: os.Getenv("INFLUXDB_DATABASE"),
	}
	var res []client.Result
	if response, err := influxClient.Query(q); err == nil {
		if response.Error() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		res = response.Results
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	c.JSON(http.StatusOK, res)
	return
}
