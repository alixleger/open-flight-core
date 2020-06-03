package flightpricesupdater

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	models "github.com/alixleger/open-flight-core/db"
	"github.com/alixleger/open-flight-core/server/handlers"
	"github.com/alixleger/open-flight-core/services/skyscanner"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	"github.com/mailgun/mailgun-go/v4"
)

// FlightPricesUpdater type
type FlightPricesUpdater struct {
	sqlDB            *gorm.DB
	skyscannerClient *skyscanner.Client
	influxdbClient   *client.Client
}

// New function is the FlightPricesUpdater constructor
func New(db *gorm.DB, skyscannerClient *skyscanner.Client, influxdbClient *client.Client) *FlightPricesUpdater {
	return &FlightPricesUpdater{sqlDB: db, skyscannerClient: skyscannerClient, influxdbClient: influxdbClient}
}

// Run flight prices updater
func (updater *FlightPricesUpdater) Run() {
	var favFlights []models.FavFlight
	updater.sqlDB.Preload("User").Preload("Flight").Find(&favFlights)

	for _, favFlight := range favFlights {
		flight := favFlight.Flight

		var departPlace models.Place
		var arrivalPlace models.Place
		updater.sqlDB.First(&departPlace, flight.DepartPlaceID)
		updater.sqlDB.First(&arrivalPlace, flight.ArrivalPlaceID)

		outboundDate := fmt.Sprintf("%d-%02d-%02d", flight.DepartDate.Year(), flight.DepartDate.Month(), flight.DepartDate.Day())
		skyscannerFlights, err := updater.skyscannerClient.GetFlights(departPlace.ExternalID, arrivalPlace.ExternalID, outboundDate)

		if err != nil {
			log.Println(err)
			continue
		}

		if len(skyscannerFlights) == 0 {
			log.Println("Skyscanner client does not return data for flight " + flight.ExternalID)
			continue
		}

		priceUpdated := false
		for _, skyscannerFlight := range skyscannerFlights {
			if skyscannerFlight.Carrier != flight.Carrier {
				continue
			}

			flight.Price = uint(skyscannerFlight.Price)
			updater.sqlDB.Save(&flight)
			priceUpdated = true
			break
		}

		if !priceUpdated {
			log.Println("Skyscanner client does not return data for flight_id " + flight.ExternalID)
			continue
		}

		if flight.Price <= favFlight.TargetPrice {
			mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_PRIVATE_KEY"))
			mg.SetAPIBase(mailgun.APIBaseEU)
			sender := os.Getenv("MAILGUN_SENDER")
			subject := "Prix cible atteint pour l'un de vos vols favoris !"
			recipient := favFlight.User.Email
			body := fmt.Sprintf("Le vol %s du %s reliant %s à %s avec la company %s a atteint votre prix cible : %d€ !", flight.ExternalID, outboundDate, departPlace.Name, arrivalPlace.Name, flight.Carrier, flight.Price)
			message := mg.NewMessage(sender, subject, body, recipient)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, id, err := mg.Send(ctx, message)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("ID: %s Resp: %s\n", id, resp)
		}

		influxClient := *updater.influxdbClient
		err = handlers.InsertFlightPrice(influxClient, favFlight.User.Email, flight.ExternalID, flight.Price)
		if err != nil {
			log.Println(err)
		}
	}
}
