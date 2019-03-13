package main

import (
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	testEmail := "tonyyanga@gmail.com"
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
