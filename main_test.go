package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_setupRouter(t *testing.T) {
	router := SetRouter()

	w := httptest.NewRecorder()
	// req1, _ := http.NewRequest("GET", "/hc", nil)
	// router.ServeHTTP(w, req1)

	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.Contains(t, w.Body.String(), "Health Check")

	// var jsonStr1 = []byte(`{"account":"account","password":"password", "email":"test123@gmail.com"}`)
	// req2, _ := http.NewRequest("POST", "/v1/users/", bytes.NewBuffer(jsonStr1))

	// router.ServeHTTP(w, req2)
	// assert.Equal(t, http.StatusOK, w.Code)

	var jsonStr2 = []byte(`{"account":"account", "password":"password"}`)
	req3, _:= http.NewRequest("POST", "/v1/users/login/", bytes.NewBuffer(jsonStr2))
	router.ServeHTTP(w, req3)
	assert.Equal(t, http.StatusOK, w.Code)
}