package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" // DB driver
	"github.com/joho/godotenv"
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

	db.AutoMigrate(&User{}, &Company{}, &Airport{}, &Flight{}, &FavFlight{})

	return db
}

// SetupTestModels function
func SetupTestModels() *gorm.DB {
	err := godotenv.Load("../../.env")

	log.Println(os.Getenv("TEST_DB_HOST"))
	log.Println(os.Getenv("TEST_DB_PORT"))
	log.Println(os.Getenv("TEST_DB_USER"))
	log.Println(os.Getenv("TEST_DB_NAME"))
	log.Println(os.Getenv("TEST_DB_PASSWORD"))

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

	db.AutoMigrate(&User{}, &Company{}, &Airport{}, &Flight{}, &FavFlight{})

	return db
}
