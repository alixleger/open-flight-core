package main

import (
	"github.com/alixleger/open-flight/back/api"
	"github.com/alixleger/open-flight/back/api/models"
	_ "github.com/joho/godotenv/autoload" // .env autoloading
)

func main() {
	server := api.New(models.SetupModels())
	server.Run()
}
