package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	gorm.Model
	Email      string      `gorm:"size:100;not null;unique_index" json:"email"`
	Password   string      `gorm:"size:100;not null;" json:"password"`
	FavFlights []FavFlight `gorm:"foreignkey:UserID" json:"fav_flights"`
}

// Hash function
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword function
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
