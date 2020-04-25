package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" // DB driver
)

// SetupModels function
func SetupModels() *gorm.DB {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD"),
		),
	)
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&User{}, &Company{}, &Flight{}, &FavFlight{})

	db.Model(&FavFlight{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&FavFlight{}).AddForeignKey("flight_id", "flights(id)", "RESTRICT", "RESTRICT")
	db.Model(&Flight{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")

	return db
}

// SetupTestModels function
func SetupTestModels() *gorm.DB {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			os.Getenv("TEST_DB_HOST"),
			os.Getenv("TEST_DB_PORT"),
			os.Getenv("TEST_DB_USER"),
			os.Getenv("TEST_DB_NAME"),
			os.Getenv("TEST_DB_PASSWORD"),
		),
	)
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&User{}, &Company{}, &Flight{}, &FavFlight{})

	db.Model(&FavFlight{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&FavFlight{}).AddForeignKey("flight_id", "flights(id)", "RESTRICT", "RESTRICT")
	db.Model(&Flight{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")

	return db
}
