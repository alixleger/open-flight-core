package skyscanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Client type
type Client struct {
	apiURL       string
	rapidAPIHost string
	rapidAPIKey  string
}

// New function is the skyscanner client constructor
func New(apiURL string, rapidAPIHost string, rapidAPIKey string) *Client {
	return &Client{
		apiURL:       apiURL,
		rapidAPIHost: rapidAPIHost,
		rapidAPIKey:  rapidAPIKey,
	}
}

func (client *Client) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", client.apiURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-host", client.rapidAPIHost)
	req.Header.Add("x-rapidapi-key", client.rapidAPIKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetPlaces return Skyscanner places from query
func (client *Client) GetPlaces(query string) ([]Place, error) {
	response, err := client.get(fmt.Sprintf("apiservices/autosuggest/v1.0/FR/EUR/fr-FR/?query=%s", query))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var structuredResponse map[string][]Place
	err = json.Unmarshal(response, &structuredResponse)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return structuredResponse["Places"], nil
}

// GetFlights return Skyscanner quotes for an itinerary
func (client *Client) GetFlights(originPlace string, destinationPlace string, outboundDate string) ([]Flight, error) {
	response, err := client.get(fmt.Sprintf("apiservices/browsequotes/v1.0/FR/EUR/fr-FR/%s/%s/%s", originPlace, destinationPlace, outboundDate))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var structuredResponse map[string]json.RawMessage
	err = json.Unmarshal(response, &structuredResponse)
	if err != nil {
		log.Println(err.Error() + " when trying to unmarshal into structuredResponse. Body: " + string(response))
		return nil, err
	}

	if val, ok := structuredResponse["message"]; ok {
		log.Println(string(val))
		return nil, errors.New(string(val))
	}

	var quotes []Quote
	var places []QuotePlace
	var carriers []Carrier

	err = json.Unmarshal(structuredResponse["Quotes"], &quotes)
	if err != nil {
		log.Println(err.Error() + " when trying to unmarshal into quotes. Body: " + string(response))
		return nil, err
	}
	err = json.Unmarshal(structuredResponse["Places"], &places)
	if err != nil {
		log.Println(err.Error() + " when trying to unmarshal into places. Body: " + string(response))
		return nil, err
	}
	err = json.Unmarshal(structuredResponse["Carriers"], &carriers)
	if err != nil {
		log.Println(err.Error() + " when trying to unmarshal into carriers. Body: " + string(response))
		return nil, err
	}

	var flights []Flight

	for _, quote := range quotes {
		var flight Flight
		flight.Price = quote.MinPrice
		carrierID := quote.OutboundLeg.CarrierIds[0]
		for _, carrier := range carriers {
			if carrier.CarrierId != carrierID {
				continue
			}
			flight.Carrier = carrier.Name
			break
		}
		for _, place := range places {
			if place.PlaceId == quote.OutboundLeg.DestinationId {
				flight.DestinationPlace = Place{
					PlaceId:     fmt.Sprintf("%s-sky", place.SkyscannerCode),
					PlaceName:   place.Name,
					CountryId:   "",
					RegionId:    "",
					CityId:      place.CityId,
					CountryName: place.CountryName,
				}
			} else if place.PlaceId == quote.OutboundLeg.OriginId {
				flight.OriginPlace = Place{
					PlaceId:     fmt.Sprintf("%s-sky", place.SkyscannerCode),
					PlaceName:   place.Name,
					CityId:      place.CityId,
					CountryName: place.CountryName,
				}
			}
		}
		flight.DepartureDate, err = time.Parse(
			time.RFC3339,
			quote.OutboundLeg.DepartureDate+"Z")
		if err != nil {
			log.Println(err.Error() + " when trying to parse DepartureDate : " + quote.OutboundLeg.DepartureDate)
			return nil, err
		}

		flights = append(flights, flight)
	}

	return flights, nil
}
