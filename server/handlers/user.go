package handlers

import (
	"net/http"

	models "github.com/alixleger/open-flight-core/db"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"gopkg.in/go-playground/validator.v9"
)

// IdentityKey of user
var IdentityKey = "id"

type loginInput struct {
	Email    string `form:"email" json:"email" binding:"required" validate:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Authenticate a user
func Authenticate(c *gin.Context) (interface{}, error) {
	var loginVals loginInput
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	email := loginVals.Email
	password := loginVals.Password

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil || models.VerifyPassword(user.Password, password) != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &user, nil
}

// Register a user
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input loginInput
	validator := validator.New()
	if err := c.ShouldBindJSON(&input); err != nil || validator.Struct(input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should have a correct email and a password fields"})
		return
	}

	// Check if email already exist
	var user models.User
	if !db.Where("email = ?", input.Email).First(&user).RecordNotFound() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This email is already taken."})
		return
	}

	hashedPassword, err := models.Hash(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot hash password"})
		return
	}

	// Create user
	user = models.User{Email: input.Email, Password: string(hashedPassword)}
	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
	return
}

// PatchUser function update user informations
func PatchUser(c *gin.Context) {
	userInterface, _ := c.Get(IdentityKey)
	user := userInterface.(*models.User)

	type inputType struct {
		Email    string `form:"email" json:"email" validate:"email"`
		Password string `form:"password" json:"password"`
	}

	var input inputType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should have a correct email and a password fields"})
		return
	}

	if input.Email != "" {
		validator := validator.New()
		if validator.Struct(input) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request should have a correct email field"})
		}
		user.Email = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := models.Hash(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Save(&user)
}
