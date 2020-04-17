package handlers

import (
	"net/http"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type companiesType struct {
	Companies []models.Company `json:"companies"`
}

// GetCompanies function
func GetCompanies(c *gin.Context) {
	var companies []models.Company
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&companies)
	var res companiesType
	res.Companies = companies

	c.JSON(http.StatusOK, res)
	return
}
