package apitests

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/valyala/fastjson"

	"github.com/stretchr/testify/assert"
)

func TestCompaniesEmpty(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}

	r := gofight.New()
	r.GET("/companies").
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)

			var p fastjson.Parser
			v, err := p.Parse(r.Body.String())
			if err != nil {
				log.Fatal(err)
			}

			assert.Len(t, v.GetArray("companies"), 0)
		})
}

func TestCompaniesNotEmpty(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}

	max := rand.Intn(100)
	for i := 0; i < max; i++ {
		createCompany("Company test " + strconv.Itoa(i))
	}

	r := gofight.New()
	r.GET("/companies").
		Run(Server.Router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)

			var p fastjson.Parser
			v, err := p.Parse(r.Body.String())

			if err != nil {
				log.Fatal(err)
			}

			assert.Len(t, v.GetArray("companies"), max)
		})
}
