package controllertests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alixleger/open-flight/back/api"
	"github.com/alixleger/open-flight/back/api/models"
	"github.com/joho/godotenv"
)

var server = api.Server{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	server = api.New(models.SetupTestModels())

	os.Exit(m.Run())
}

func refreshBookTable() error {
	err := server.DB.DropTableIfExists(&models.Book{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Book{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table(s)")
	return nil
}

func performRequest(method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	return w
}
