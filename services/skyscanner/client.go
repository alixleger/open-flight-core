package skyscanner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (client *Client) get(endpoint string) []byte {
	url := fmt.Sprintf("%s/%s", client.apiURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-rapidapi-host", client.rapidAPIHost)
	req.Header.Add("x-rapidapi-key", client.rapidAPIKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

// Place ressource
type Place struct {
	PlaceId     string
	PlaceName   string
	CountryId   string
	RegionId    string
	CityId      string
	CountryName string
}

// GetPlaces return Skyscanner places from query
func (client *Client) GetPlaces(query string) []Place {
	response := client.get(fmt.Sprintf("/apiservices/autosuggest/v1.0/FR/EUR/fr-FR/?query=%s", query))
	var structuredResponse map[string][]Place
	err := json.Unmarshal(response, &structuredResponse)
	if err != nil {
		log.Fatal(err)
	}

	return structuredResponse["Places"]
}

// GetQuotes return Skyscanner quotes for an itinery
func (client *Client) GetQuotes(originPlace string, destinationPlace string, outboundDate string) []byte {
	return client.get(fmt.Sprintf("apiservices/browsequotes/v1.0/FR/EUR/fr-FR/%s/%s/%s", originPlace, destinationPlace, outboundDate))
}
