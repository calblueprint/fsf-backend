package main

/*
 * NOT WORKING on local Macbooks
 * only works when run from root@fsfmobile0p.fsf.org
 * To test, run
 * TCUSERNAME="your_tc_username" TCPASSWORD="your_tc_password" SITEKEY="your_site_key" ADMINAPIKEY="your_admin_api_key" go test
 */

import (
	"os"
	"testing"
)

// test getAPIKey
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
		t.Log(apiKey)
		t.Log(contactId)
	}
}
