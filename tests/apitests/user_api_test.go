package apitests

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type loginInput struct {
	Email    string
	Password string
}

func TestRegisterOK(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	login := loginInput{
		Email:    "test@test.fr",
		Password: "test",
	}

	body, err := json.Marshal(login)

	if err != nil {
		log.Fatal(err)
	}

	w := performRequest("POST", "/register", body, "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterBadRequest(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	login := loginInput{
		Email:    "test",
		Password: "test",
	}

	body, err := json.Marshal(login)

	if err != nil {
		log.Fatal(err)
	}

	w := performRequest("POST", "/register", body, "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUnauthorized(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	login := loginInput{
		Email:    "test@test.fr",
		Password: "test",
	}

	createUser(login.Email, login.Password)

	body, err := json.Marshal(login)

	if err != nil {
		log.Fatal(err)
	}

	w := performRequest("POST", "/register", body, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
