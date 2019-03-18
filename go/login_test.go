package main

/*
 * NOT WORKING on local Macs
 * only works when run from root@fsfmobile0p.fsf.org
 */

import (
	"os"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	testEmail := "tonyyanga@gmail.com"
	tcUsername = os.Getenv("TCUSERNAME")
	tcPassword = os.Getenv("TCPASSWORD")
	siteKey = os.Getenv("SITEKEY")
	adminAPIKey = os.Getenv("ADMINAPIKEY")
	apiKey, contactId, err := getAPIKey(testEmail)
	if err != nil {
		t.Error(err.Error())
	} else {
		// t.Log("Tony's apikey is %s", apiKey)
		// t.Log("Tony's contactId is %s", contactId)
		t.Log(apiKey)
		t.Log(contactId)
	}
}
