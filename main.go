package main

import (
	"log"
	"os"

	"github.com/alixleger/open-flight-core/db"
	"github.com/alixleger/open-flight-core/server"
	"github.com/alixleger/open-flight-core/services/skyscanner"
	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	sqlDB := db.SetupModels()
	skyscannerClient := skyscanner.New(os.Getenv("SKYSCANNER_API_URL"), os.Getenv("SKYSCANNER_API_HOST"), os.Getenv("SKYSCANNER_API_KEY"))
	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("INFLUXDB_ADDR"),
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	})

	if err != nil {
		log.Fatal(err)
	}
	defer influxdbClient.Close()

	server := server.New(sqlDB, skyscannerClient, &influxdbClient)
	server.Run()
}
