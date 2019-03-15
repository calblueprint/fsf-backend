package main

import (
	"encoding/xml"
	"log"
	"net/http"
)

// This function accepts a service token from the client, who obtains it from
// CAS. It interacts with CAS server to validate the token.
//
// @param token: service token
// @return:
//   bool: whether auth is successful
//   string: id string to retrieve user auth info from civicrm
//   error: if err != nil, there is a server error
func validateToken(token string) (bool, string, error) {
	c := &http.Client{}

	requestURLPrefix := "https://cas.fsf.org/serviceValidate?service=https://crmserver3d.fsf.org/associate/account&ticket="

	requestURL := requestURLPrefix + token

	resp, err := c.Get(requestURL)

	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return false, "", err
	}

	if resp.StatusCode >= 400 {
		return false, "", nil
	}

	// resp StatusCode should be 200, parse response for identity

	dec := xml.NewDecoder(resp.Body)

	var r struct {
		XMLName     xml.Name `xml:"serviceResponse"`
		AuthSuccess struct {
			XMLName    xml.Name `xml:"authenticationSuccess"`
			User       string   `xml:"user"`
			Attributes struct {
				XMLName xml.Name `xml:"attributes"`
				Email   string   `xml:"mail"`
			} `xml:"attributes"`
		} `xml:"authenticationSuccess"`
	}

	if err := dec.Decode(&r); err != nil {
		log.Fatal("Failure parsing CAS validate response")
		return false, "", nil
	}

	return true, r.AuthSuccess.Attributes.Email, nil
}
