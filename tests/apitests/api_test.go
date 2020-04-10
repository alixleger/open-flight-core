package apitests

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alixleger/open-flight/back/api"
	"github.com/alixleger/open-flight/back/api/models"
)

var server = api.Server{}

func TestMain(m *testing.M) {
	server = api.New(models.SetupTestModels())
	os.Exit(m.Run())
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func createUser(email string, password string) {
	hashedPassword, err := models.Hash(password)
	if err != nil {
		log.Fatal(err)
	}
	user := models.User{Email: email, Password: string(hashedPassword)}
	server.DB.Create(&user)
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
	server.Router.ServeHTTP(w, req)

	return w
}
