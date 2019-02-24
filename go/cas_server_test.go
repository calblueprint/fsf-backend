package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Test /payment/register endpoint
func TestHandleRegisterCC(t *testing.T) {
	tcUsername = os.Getenv("TCUSERNAME")
	tcPassword = os.Getenv("TCPASSWORD")
	prePayload := []byte(`{"name": "Delete Me After Test", "cc": "4111111111111111", "exp": "0223", "zip": "90000", "demo": "y"}`)
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
