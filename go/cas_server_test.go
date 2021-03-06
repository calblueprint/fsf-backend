package main

/*
 * NOT WORKING on local Macbooks
 * only works when run from root@fsfmobile0p.fsf.org
 * To test, run
 * TCUSERNAME="your_tc_username" TCPASSWORD="your_tc_password" go test
 */

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type PayStruct struct {
	Name   string `json:"name"`
	Cc     string `json:"cc"`
	Cvv    string `json:"cvv"`
	Exp    string `json:"exp"`
	Amount string `json:"amount"`
	Email  string `json:"email"`
	Apikey string `json:"apikey"`
}

type RepeatPayStruct struct {
	BillingID string `json:"billingid"`
	Amount    string `json:"amount"`
	Email     string `json:"email"`
	Apikey    string `json:"apikey"`
}

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

// test /payment/pay endpoint
func TestHandlePayment(t *testing.T) {
	tcUsername = os.Getenv("TCUSERNAME")
	tcPassword = os.Getenv("TCPASSWORD")
	siteKey = os.Getenv("SITEKEY")
	adminAPIKey = os.Getenv("ADMINAPIKEY")

	payStruct := PayStruct{Name: "John Smith", Cc: "4111111111111111", Cvv: "123", Exp: "0404", Amount: "110", Email: "tonyyanga@gmail.com", Apikey: adminAPIKey}
	prePayload, err := json.Marshal(payStruct)

	payLoad := bytes.NewBuffer(prePayload)

	request, err := http.NewRequest("POST", "/payment/pay", payLoad)

	if err == nil {
		handler := http.HandlerFunc(handlePayment)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Errorf("/payment/pay FAILED")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestPaymentFromBillingId(t *testing.T) {
	tcUsername = os.Getenv("TCUSERNAME")
	tcPassword = os.Getenv("TCPASSWORD")
	siteKey = os.Getenv("SITEKEY")
	adminAPIKey = os.Getenv("ADMINAPIKEY")

	repeatPayStruct := RepeatPayStruct{BillingID: "Q50K8A", Amount: "5315", Email: "tonyyanga@gmail.com", Apikey: adminAPIKey}
	prePayload, err := json.Marshal(repeatPayStruct)
	if err != nil {
		t.Errorf("error encoding json")
	}

	payLoad := bytes.NewBuffer(prePayload)
	request, err := http.NewRequest("POST", "/payment/repeat_pay", payLoad)

	if err == nil {
		handler := http.HandlerFunc(handleRepeatPayment)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Errorf("/payment/repeat_pay FAILED")
		}
	} else {
		t.Errorf(err.Error())
	}
}
