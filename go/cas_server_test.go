package main

/*
 * NOT WORKING on local Macbooks
 * only works when run from root@fsfmobile0p.fsf.org
 * To test, run
 * TCUSERNAME="your_tc_username" TCPASSWORD="your_tc_password" go test
 */

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
			// t.Errorf("%i", response.Code)
		}
	} else {
		t.Errorf(err.Error())
	}
}

// test /payment/pay endpoint
func TestHandlePayment(t *testing.T) {
	tcUsername = os.Getenv("TCUSERNAME")
	tcPassword = os.Getenv("TCPASSWORD")
	siteKey = os.Getenv("SITEKEY")
	adminAPIKey = os.Getenv("ADMINAPIKEY")

	prePayload := []byte(`{"name": "John Smith", "cc": "4111111111111111", "exp": "0404", "email": "test@test.com", "apikey": ""}`)
	payLoad := bytes.NewBuffer(prePayload)

	request, err := http.NewRequest("POST", "/payment/pay", payLoad)

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
