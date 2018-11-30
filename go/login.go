package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var siteKey, adminAPIKey string

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

// A helper function to query CiviCRM
// @param
//   v: encoded CiviCRM REST query
//   dest: an object where we store the decoded json object
func queryCiviCRM(v url.Values, dest interface{}) error {
	c := &http.Client{}
	requestURL := "https://crmserver3d.fsf.org/sites/all/modules/civicrm/extern/rest.php"

	resp, err := c.PostForm(requestURL, v)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Bad status code: %v", string(resp.StatusCode))
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(dest)
	if err != nil {
		return fmt.Errorf("Bad response")
	}

	return nil
}

// Retrieve the API key of the user identified by certain id string
// @param id: id string from CAS (currently user's email)
// @return:
//   a string that is the API key
//   a string that is the contact id in CiviCRM
//   an error if any error occurs
func getAPIKey(id string) (string, string, error) {
	/** Because of the design of CiviCRM, API Key is only shown when we do an update
	  i.e. a create with contact_id specified.

	  We first query for the contact id that matches this identity, then run an
	  update to retrieve the API key.

	  If current API key is not set, we set the API key and return to the user.
	*/

	// Query for contact id
	var idQuery struct {
		Sequential int    `json:"sequential"`
		Return     string `json:"return"`
		Email      string `json:"email"`
	}
	idQuery.Sequential = 1
	idQuery.Email = id
	idQuery.Return = "id"

	idQueryJson, err := json.Marshal(idQuery)
	if err != nil {
		log.Fatal("Error constructing query json for civicrm")
		return "", "", err
	}

	v := &url.Values{}
	v.Add("entity", "Contact")
	v.Add("action", "get")
	v.Add("api_key", adminAPIKey)
	v.Add("key", siteKey)
	v.Add("json", string(idQueryJson))

	var idQueryResp struct {
		Error  int `json:"is_error"`
		Values []struct {
			Id string `json:"id"`
		} `json:"values"`
	}

	if err = queryCiviCRM(*v, &idQueryResp); err != nil || idQueryResp.Error != 0 || len(idQueryResp.Values) != 1 {
		return "", "", fmt.Errorf("Bad response")
	}

	contactId := idQueryResp.Values[0].Id
	log.Println("Contact id is:" + contactId)

	// Use an update query to check for API key
	updateQueryJson := `{"id":"` + contactId + `"}`

	v.Set("action", "create")
	v.Set("json", updateQueryJson)

	var updateQueryResp struct {
		Error  int `json:"is_error"`
		Values map[string]struct {
			APIKey string `json:"api_key"`
		} `json:"values"`
	}

	if err = queryCiviCRM(*v, &updateQueryResp); err != nil || updateQueryResp.Error != 0 {
		return "", "", fmt.Errorf("Bad response")
	}

	if updateQueryResp.Values[contactId].APIKey != "" {
		return updateQueryResp.Values[contactId].APIKey, contactId, nil
	}

	log.Println("Found API key:" + updateQueryResp.Values[contactId].APIKey)
	// API key is not set, need to update API key and return
	newAPIKey := getRandomKey()

	updateQueryJson = `{"id":"` + contactId + `", "api_key": "` + newAPIKey + `"}`
	v.Set("json", updateQueryJson)

	if err = queryCiviCRM(*v, &updateQueryResp); err != nil || updateQueryResp.Error != 0 {
		return "", "", fmt.Errorf("Bad response")
	}

	log.Println("Set API key:" + updateQueryResp.Values[contactId].APIKey)

	return newAPIKey, contactId, nil
}
