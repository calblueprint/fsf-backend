package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

// Test /payment/register endpoint
func TestHandleRegisterCC(t *testing.T) { // it needs to begin with Test
	tcUsername = os.Getenv("USERNAME")
	tcPassword = os.Getenv("PW")
	prePayload := []byte(`{"name": "Delete Me After Test", "cc": "4111111111111111", "exp": "0223", "zip": "90000"}`)
	payLoad := bytes.NewBuffer(prePayload)

	request, err := http.NewRequest("POST", "/payment/register", payLoad)

	if err == nil {
		handler := http.HandlerFunc(handleRegisterCC)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Errorf("/payment/register FAILED")
		}
	} else {
		t.Errorf(err.Error())
	}
}

