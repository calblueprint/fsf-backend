package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

// global variables to be set by users passing in env variables
var siteKey, adminAPIKey string

/*
A helper function to query CiviCRM
@param
  v: encoded CiviCRM REST query
  dest: an object where we store the decoded json object
*/
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
		log.Println(err)
		return fmt.Errorf("Bad response")
	}

	return nil
}

/*
	Retrieve the API key of the user identified by certain id string
	@param id: id string from CAS (currently user's email)
	@return:
	a string that is the API key
	a string that is the contact id in CiviCRM
	an error if any error occurs
*/
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
		log.Log("Error constructing query json for civicrm")
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
		Error  int                 `json:"is_error"`
		Values map[string]UserInfo `json:"values"`
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

/*
	A helper function to validate a user apiKey that we get from the client
	@param:
	apiKey: apiKey stored in frontend
	contactId: contactId of user
*/
func validateAPIKeyForUpdateRequests(apiKey string, contactId string) bool {
	idQueryJSON := `{"id":"` + contactId + `"}`
	v := &url.Values{}
	v.Add("entity", "Contact")
	v.Add("api_key", adminAPIKey)
	v.Add("key", siteKey)
	v.Add("action", "create")
	v.Add("json", idQueryJSON)

	var updateQueryResp struct {
		Error  int                 `json:"is_error"`
		Values map[string]UserInfo `json:"values"`
	}
	err := queryCiviCRM(*v, &updateQueryResp)

	if err != nil || updateQueryResp.Error != 0 {
		return false
	}

	if updateQueryResp.Values[contactId].APIKey != "" {
		return updateQueryResp.Values[contactId].APIKey == apiKey
	}

	return false
}

/*
	A helper function to get a user's info
	@param:
	apiKey: apiKey stored in frontend
	contactId: contactId of user
	@return:
	userinfos struct, error
*/
func getUserInfo(apiKey string, contactId string) (*UserInfo, error) {
	/** Because of the design of CiviCRM, API Key is only shown when we do an update
	  i.e. a create with contact_id specified.

	  We first query for the contact id that matches this identity, then run an
	  update to retrieve the API key.

	  If current API key is not set, we set the API key and return to the user.
	*/
	valid := validateAPIKeyForUpdateRequests(apiKey, contactId)
	if !valid {
		return nil, fmt.Errorf("Bad api key")
	}
	idJSON := `{"id":"` + contactId + `"}`
	v := &url.Values{}
	v.Add("entity", "Contact")
	v.Add("action", "get")
	v.Add("api_key", adminAPIKey)
	v.Add("key", siteKey)
	v.Add("json", idJSON)

	var getQueryResp struct {
		Error  int                 `json:"is_error"`
		Values map[string]UserInfo `json:"values"`
	}
	err := queryCiviCRM(*v, &getQueryResp)
	if err != nil || getQueryResp.Error != 0 {
		log.Println(err)
		log.Println(getQueryResp.Error)
		return nil, fmt.Errorf("Bad response")
	}
	var userInfo = getQueryResp.Values[contactId]
	return &userInfo, nil
}

// records transaction in CiviCRM
/*
	A helper function to gerecord a transaction in CiviCRM.
	@param:
	userEmail, userAPIKey, transID, amount
	@return:
	error
	The function can be edited to store whatever information we want from a transaction
*/
func recordTransactionInCiviCRM(userEmail string, userAPIKey string, transID string, amount string) error {
	// record this transaction in CiviCRM
	civiCRMAPIKey, userContactId, err := getAPIKey(userEmail)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error retrieving contact info from CiviCRM")
	}

	// validate API key received form the client with API key in CiviCRM
	if !reflect.DeepEqual(userAPIKey, civiCRMAPIKey) {
		return errors.New("authentication failed - api keys do not match")
	}

	// transactionInfo struct to put into CiviCRM
	var transactionInfo struct {
		FinancialTypeId string `json:"financial_type_id"`
		TotalAmount     string `json:"total_amount"`
		ContactId       string `json:"contact_id"`
		TrxnId          string `json:"trxn_id"`
		// make edit here to store more/different information in CiviCRM
	}

	// this particular transaction is a donation
	transactionInfo.FinancialTypeId = "Donation"
	transactionInfo.TotalAmount = amount
	transactionInfo.ContactId = userContactId
	transactionInfo.TrxnId = transID

	infoJson, err := json.Marshal(transactionInfo)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error constructing infoJson for civicrm from transactionInfo")
	}

	v := &url.Values{}
	v.Add("entity", "Contribution")
	v.Add("action", "create")
	v.Add("api_key", adminAPIKey)
	v.Add("key", siteKey)
	v.Add("json", string(infoJson))

	var infoPutResp struct {
		Error int `json:"is_error"`
	}

	if err = queryCiviCRM(*v, &infoPutResp); err != nil {
		return err
	} else if infoPutResp.Error != 0 {
		return errors.New("error querying CiviCRM")
	}

	return nil
}

// Useful Structs
// Not all of these fields will be populated
type UserInfo struct {
	ID                           string `json:"id"`
	ContactType                  string `json:"contact_type"`
	ContactSubType               string `json:"contact_sib_type"`
	DoNotEmail                   string `json:"do_not_email"`
	DoNotPhone                   string `json:"do_not_phone"`
	DoNotMail                    string `json:"do_not_mail"`
	DoNotSms                     string `json:"do_not_sms"`
	DoNotTrade                   string `json:"do_not_trade"`
	IsOptOut                     string `json:"is_opt_out"`
	LegalIdentifier              string `json:"legal_identifier"`
	ExternalIdentifier           string `json:"external_identifier"`
	SortName                     string `json:"sort_name"`
	DisplayName                  string `json:"display_name"`
	NickName                     string `json:"nick_name"`
	LegalName                    string `json:"legal_name"`
	ImageUrl                     string `json:"image_URL"`
	PreferredCommunicationMethod string `json:"preferred_communication_method"`
	PreferredLanguage            string `json:"preferred_language"`
	PreferredMailFormat          string `json:"preferred_mail_format"`
	Hash                         string `json:"hash"`
	APIKey                       string `json:"api_key"`
	FirstName                    string `json:"first_name"`
	MiddleName                   string `json:"middle_name"`
	LastName                     string `json:"last_name"`
	PrefixId                     string `json:"prefix_id"`
	SuffixId                     string `json:"suffix_id"`
	FormalTitle                  string `json:"formal_title"`
	CommunicationStyleId         string `json:"communication_style_id"`
	JobTitle                     string `json:"job_title"`
	GenderId                     string `json:"gender_id"`
	BirthDate                    string `json:"birth_date"`
	IsDeceased                   string `json:"is_deceased"`
	DeceasedDate                 string `json:"deceased_date"`
	HouseholdName                string `json:"household_name"`
	OrganizationName             string `json:"organization_name"`
	SicCode                      string `json:"sic_code"`
	ContactIsDeleted             string `json:"contact_is_deleted"`
	CurrentEmployer              string `json:"current_employer"`
	AddressId                    string `json:"address_id"`
	StreetAddress                string `json:"street_address"`
	SupplementalAddress1         string `json:"supplemental_address_1"`
	SupplementalAddress2         string `json:"supplemental_address_2"`
	SupplementalAddress3         string `json:"upplemental_address_3"`
	City                         string `json:"city"`
	PostalCodeSuffix             string `json:"postal_code_suffix"`
	PostalCode                   string `json:"postal_code"`
	GeoCode1                     string `json:"geo_code_1"`
	GeoCode2                     string `json:"geo_code_2"`
	StateProvinceId              string `json:"state_province_id"`
	CountryId                    string `json:"country_id"`
	PhoneId                      string `json:"phone_id"`
	PhoneTypeId                  string `json:phone_type_id`
	Phone                        string `json:"phone"`
	EmailId                      string `json:"email_id"`
	Email                        string `json:"email"`
	OnHold                       string `json:"on_hold"`
	ImId                         string `json:"im_id"`
	ProviderId                   string `json:"provider_id"`
	Im                           string `json:"im"`
	WorldRegionId                string `json:"worldregion_id"`
	worldRegion                  string `json:"world_region"`
	languages                    string `json:"languages"`
	IndividualPrefix             string `json:individual_prefix"`
	IndividualSuffix             string `json:"individual_suffix"`
	CommunicationStyle           string `json:"communication_style"`
	Gender                       string `json:"gender"`
	StateProvinceName            string `json:"state_province_name"`
	StateProvince                string `json:"state_province"`
	Country                      string `json:"country"`
}
