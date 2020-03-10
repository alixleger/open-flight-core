package controllertests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBooks(t *testing.T) {
	err := refreshBookTable()
	if err != nil {
		log.Fatal(err)
	}

	w := performRequest("GET", "/books")
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]string
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	var errMessage string
	if err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			errMessage = fmt.Sprintf("UnmarshalTypeError %v - %v - %v\n", ute.Value, ute.Type, ute.Offset)
		} else {
			errMessage = fmt.Sprintln("Other error:", err)
		}
	}
	assert.Nil(t, err, errMessage)

	value, exists := response["data"]
	assert.True(t, exists)
	assert.Empty(t, value)

	// TODO: Insert mock books and test again
}
