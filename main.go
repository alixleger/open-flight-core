package main

import (
	"github.com/alixleger/open-flight-core/db"
	"github.com/alixleger/open-flight-core/server"
	_ "github.com/joho/godotenv/autoload" // .env autoloading
)

func main() {
	server := server.New(db.SetupModels())
	server.Run()
}
