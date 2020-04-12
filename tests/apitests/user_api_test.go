package apitests

import (
	"log"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"

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

	r := gofight.New()
	r.POST("/register").
		SetJSON(gofight.D{
			"email":    login.Email,
			"password": login.Password,
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
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

	r := gofight.New()
	r.POST("/register").
		SetJSON(gofight.D{
			"email":    login.Email,
			"password": login.Password,
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
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

	r := gofight.New()
	r.POST("/register").
		SetJSON(gofight.D{
			"email":    login.Email,
			"password": login.Password,
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusUnauthorized, r.Code)
		})
}

func TestUserPatchOK(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	r := gofight.New()

	r.PATCH("/auth/user").
		SetHeader(gofight.H{
			"Authorization": "Bearer " + createUser("test@test.fr", "test"),
		}).
		SetJSON(gofight.D{
			"email":    "test@test.fr",
			"password": "testtest",
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	r.POST("/login").
		SetJSON(gofight.D{
			"email":    "test@test.fr",
			"password": "testtest",
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUserPatchUnauthorized(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	createUser("test@test.fr", "test")

	r := gofight.New()
	r.PATCH("/auth/user").
		SetHeader(gofight.H{
			"Authorization": "Bearer " + "test",
		}).
		SetJSON(gofight.D{
			"email":    "test@test.fr",
			"password": "testtest",
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusUnauthorized, r.Code)
		})
}

func TestUserPatchBadRequest(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	r := gofight.New()
	r.PATCH("/auth/user").
		SetHeader(gofight.H{
			"Authorization": "Bearer " + createUser("test@test.fr", "test"),
		}).
		SetJSON(gofight.D{
			"email":    "test",
			"password": "",
		}).
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
