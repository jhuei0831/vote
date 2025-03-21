package main

import (
	"bytes"
	"encoding/json"
	"log"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var token string

func Test_setupRouter(t *testing.T) {
	router := SetRouter()

	w := httptest.NewRecorder()
	// req1, _ := http.NewRequest("GET", "/hc", nil)
	// router.ServeHTTP(w, req1)

	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.Contains(t, w.Body.String(), "health check: PORT 9443")

	var jsonStr2 = []byte(`{"account":"creator", "password":"password"}`)
	req3, _:= http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(jsonStr2))
	router.ServeHTTP(w, req3)
	// get token
	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	token = response["data"].(map[string]any)["token"].(string)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateUser(t *testing.T) {
	router := SetRouter()
	w := httptest.NewRecorder()
	var jsonStr1 = []byte(`{"account":"account","password":"password", "email":"test123@gmail.com"}`)
	req2, err := http.NewRequest("POST", "/v1/user/create", bytes.NewBuffer(jsonStr1))
	if err != nil {
		log.Fatal(err)
	}
	req2.Header.Set("Content-Type", "application/json")
	// add token to header
	req2.Header.Set("Authorization", "Bearer " + token)

	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Body.String())
}