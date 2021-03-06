package apitests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	models "github.com/alixleger/open-flight-core/db"
	api "github.com/alixleger/open-flight-core/server"
	"github.com/alixleger/open-flight-core/services/skyscanner"
)

var Server *api.Server

func TestMain(m *testing.M) {
	Server = api.New(models.SetupTestModels(), skyscanner.New("", "", ""), nil)
	os.Exit(m.Run())
}

func refreshUserTable() error {
	err := Server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = Server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func createUser(email string, password string) string {
	hashedPassword, err := models.Hash(password)
	if err != nil {
		log.Fatal(err)
	}
	user := models.User{Email: email, Password: string(hashedPassword)}
	Server.DB.Create(&user)

	body, err := json.Marshal(struct {
		Email    string
		Password string
	}{email, password})

	if err != nil {
		log.Fatal(err)
	}

	w := performRequest("POST", "/login", body, "")
	var response map[string]json.RawMessage
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil || response["token"] == nil {
		log.Fatal(err)
	}

	return strings.Replace(string(response["token"]), "\"", "", -1)
}

func performRequest(method, path string, jsonBody []byte, token string) *httptest.ResponseRecorder {
	var body io.Reader = nil

	if jsonBody != nil {
		body = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	Server.Router.ServeHTTP(w, req)

	return w
}
