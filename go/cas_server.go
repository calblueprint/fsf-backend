package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

var tcUsername, tcPassword string

// handles a request for a repeated payment, i.e. using billing id
// accepts a body like the following
// requires a POST request with json payload with the following format
// {
//  "billingid": "slvkdfjasdoihgjosa",
//  "amount": "110"
// }
//
// when success, returns a json like the following:
// {
//   "transid": "a transaction id from TrustCommerce",
//   "status": "status of transaction { approved, declined, baddata, error }"
//   "authcode": "auth code for the transaction"
// }
func handleRepeatPayment(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var ccInfo struct {
		BillingID string `json:"billingid"`
		Amount    string `json:"amount"`
	}

	if err := dec.Decode(&ccInfo); err != nil {
		writeError(w, "Cannot parse request body correctly")
		return
	}

	mgr := NewTransactionMgr(tcUsername, tcPassword)
	saleResp, err := mgr.createSaleFromBillingID(ccInfo.BillingID, ccInfo.Amount)

	if err != nil {
		log.Println(err.Error())
		writeError(w, "Payment failed")
		return
	}

	// write billing id back in response
	enc := json.NewEncoder(w)

	// see TCSaleResp struct for json response struct
	enc.Encode(saleResp)
}

// handles a request for a single payment
// accepts a body like the following
// requires a POST request with json payload with the following format
// {
//  "name": "John Smith",
//  "cc": "4111111111111111",
//  "exp": "0404",
//  "amount": "110"
// }
//
// when success, returns a json like the following:
// {
//   "transid": "a transaction id from TrustCommerce",
//   "status": "status of transaction { approved, declined, baddata, error }"
//   "authcode": "auth code for the transaction"
// }
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func handlePayment(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var ccInfo struct {
		Name   string `json:"name"`
		Cc     string `json:"cc"`
		Exp    string `json:"exp"`
		Amount string `json:"amount"`
		Email  string `json:"email"`
		ApiKey string `json:"apikey"`
	}

	if err := dec.Decode(&ccInfo); err != nil {
		writeError(w, "Cannot parse request body correctly")
		return
	}
	// mgr := NewTransactionMgr(tcUsername, tcPassword)
	// saleResp, err := mgr.createSaleFromCC(ccInfo.Name, ccInfo.Cc, ccInfo.Exp, ccInfo.Amount)
	// TODO: DUMMY RES FOR TESTING REMOVE LATER
	var saleResp struct {
		TransID  string `json:"transid"`
		Status   string `json:"status"`
		AuthCode string `json:"authcode"`
	}

	saleResp.TransID = "hello"
	saleResp.Status = "approved"
	saleResp.AuthCode = "world"

	/*
		type TCSaleResp struct {
			TransID  string `json:"transid"`
			Status   string `json:"status"`
			AuthCode string `json:"authcode"`
		}
	*/

	log.Printf("STEP 1")

	if saleResp.Status != "approved" {
		// do something if transaction is not approved
		writeError(w, "Transaction not successfully approved")
		return
	} else {
		// record this transaction in CiviCRM
		userEmail := ccInfo.Email
		userAPIKey := ccInfo.ApiKey
		civiCRMAPIKey, userContactId, err := getAPIKey(userEmail)
		if err != nil {
			writeError(w, "error retrieving contact info from CiviCRM")
			return
		} else if !reflect.DeepEqual(userAPIKey, civiCRMAPIKey) {
			writeError(w, "authentication failed - api keys do not match")
			return
		}
		transID := saleResp.TransID

		log.Printf("STEP 2")
		// Need to pass in API key + contact_id
		/*
					ccInfo struct {
					Name   string `json:"name"`
					Cc     string `json:"cc"`
					Exp    string `json:"exp"`
			*		Amount string `json:"amount"`
			++	Email string `json:"email"`
			++	ApiKey string `json:"apikey"`
				}
		*/
		var info struct {
			FinancialTypeId string `json:"financial_type_id"`
			TotalAmount     string `json:"total_amount"`
			ContactId       string `json:"contact_id"`
			TrxnId          string `json:"trxn_id"`
		}

		info.FinancialTypeId = "Donation"
		info.TotalAmount = ccInfo.Amount
		info.ContactId = userContactId
		info.TrxnId = transID

		infoJson, err := json.Marshal(info)
		if err != nil {
			log.Fatal("Error constructing info json for civicrm")
			return
		}

		v := &url.Values{}
		v.Add("entity", "Contribution")
		v.Add("action", "create")
		v.Add("api_key", adminAPIKey)
		v.Add("key", siteKey)
		v.Add("json", string(infoJson))

		var infoPutResp struct {
			Error int `json:"is_error"`
			/*Values []struct {
				Id string `json:"id"`
			} `json:"values"`*/
		}

		if err = queryCiviCRM(*v, &infoPutResp); err != nil || infoPutResp.Error != 0 {
			log.Printf("REQUEST: %v", v.Get("json"))
			log.Printf("Bad response: %v", err)
		}

		log.Printf("STEP 3")
		// relevantInfo := {"financial_type_id":"","total_amount":"","contact_id":"user_contact_id"}
		/*
			firstPart := "https://crmserver3d.fsf.org/sites/all/modules/civicrm/extern/rest.php?entity=Contribution&action=create&api_key="
			secondPart := "&key="
			thirdPart := "&json="
			url := firstPart + adminAPIKey + secondPart + siteKey + thirdPart + string(infoJson)
		*/
		// url := "https://crmserver3d.fsf.org/sites/all/modules/civicrm/extern/rest.php?entity=Contribution&action=create&api_key=userkey&key=sitekey&json=" + relevantInfo
		// siteKey, adminAPIKey
		/*
			request, err := http.NewRequest("POST", url, nil)
			if err != nil {
				writeError(w, "error storing transaction info to CiviCRM")
				return
			}
		*/

	}

	// TODO: UNCOMMENT WITH LINE 98
	// if err != nil {
	// 	log.Println(err.Error())
	// 	writeError(w, "Payment failed")
	// 	return
	// }

	// write billing id back in response
	enc := json.NewEncoder(w)

	// see TCSaleResp struct for json response struct
	enc.Encode(saleResp)
}

// handles a request to store credit card info for repeating payments
// accepts a body like the following
// requires a POST request with json payload with the following format
// {
//  "name": "John Smith",
//  "cc": "4111111111111111",
//  "exp": "0404",
//  "zip": "90000"
// }
//
// when success, returns a json like the following:
// {"billingid": "a billing id from TrustCommerce"}
func handleRegisterCC(w http.ResponseWriter, req *http.Request) {
	/** requires a POST request with json payload with the following format
	  *
	  * {
	      "name": "John Smith",
	      "cc": "4111111111111111",
	      "exp": "0404",
	      "zip": "90000"
	    }
	*/
	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var ccInfo struct {
		Name string `json:"name"`
		Cc   string `json:"cc"`
		Exp  string `json:"exp"`
		Zip  string `json:"zip"`
	}

	if err := dec.Decode(&ccInfo); err != nil {
		writeError(w, "Cannot parse request body correctly")
		return
	}

	mgr := NewTransactionMgr(tcUsername, tcPassword)
	billingId, err := mgr.createBillingId(ccInfo.Name, ccInfo.Cc, ccInfo.Exp, ccInfo.Zip)

	if err != nil {
		log.Println(err.Error())
		writeError(w, "Card registration failed")
		return
	}

	// write billing id back in response
	enc := json.NewEncoder(w)

	var resp struct {
		Id string `json:"billingid"`
	}

	resp.Id = billingId

	enc.Encode(resp)
}

// handles a login request
// accepts JSON POST request with the following body:
//   {"st": "service token"}
//
// returns either a HTTP error or a json response like:
//   {
//     "key": "api key for CiviCRM",
//     "id": "contact id for CiviCRM",
//     "email": "email address"
//   }
//
// at this time, the id to talk to CiviCRM is the users' email address
func handleLogin(w http.ResponseWriter, req *http.Request) {
	/** requires a POST request with json payload with the following format
	 *
	 *  {"st": "service token"}
	 */
	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var t struct {
		ST string `json:"st"`
	}

	if err := dec.Decode(&t); err != nil {
		writeError(w, "Cannot parse request body correctly")
		return
	}

	serviceToken := t.ST // service token we get from the request

	result, id, err := validateToken(serviceToken) // id here is the email from CAS

	if err != nil {
		writeError(w, "Server error when validating token")
		return
	} else if !result {
		writeAccessDenied(w, "Ticket authentication failed")
		return
	}

	// authentication success
	APIKey, contactId, err := getAPIKey(id)
	if err != nil {
		log.Println(err.Error())
		writeError(w, "Server error when interacting with CiviCRM")
		return
	}

	// write APIKey back in response
	enc := json.NewEncoder(w)

	var key struct {
		Key   string `json:"key"`
		ID    string `json:"id"` // contact id from CiviCRM
		Email string `json:"email"`
	}

	key.Key = APIKey
	key.ID = contactId
	key.Email = id

	enc.Encode(key)
}

// handles a userInfo request
// accepts JSON GET request
// returns either a HTTP error or a json response like:
//   {
//     "firstname": "user first name",
//     "lastname": "user last name",
//     "address": "user address",
//     "email": "email address"
//   }
func getUserInformation(w http.ResponseWriter, req *http.Request) {
	/** requires a POST request
	 */
	if req.Method != "GET" {
		writeError(w, "Only GET requests are supported")
		return
	}
	// write APIKey back in response
	enc := json.NewEncoder(w)

	var userInfo struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"` // contact id from CiviCRM
		Address   string `json:"address"`
		Email     string `json:"email"`
	}

	userInfo.FirstName = "Mukil"
	userInfo.LastName = "Loganathan"
	userInfo.Address = "3664 Cody Court"
	userInfo.Email = "mukil.loganathan@gmail.com"

	enc.Encode(userInfo)
}

func main() {
	// CAS mobile login server
	addrPtr := flag.String("addr", "0.0.0.0:8080", "address to listen")
	siteKeyPtr := flag.String("sitekey", "", "provide site key of CiviCRM")
	adminAPIKeyPtr := flag.String("adminkey", "", "provide the admin API key for CiviCRM")
	tcuserPtr := flag.String("tcuser", "", "username for Trust Commerce")
	tcpasswdPtr := flag.String("tcpasswd", "", "password for Trust Commerce")

	flag.Parse()

	siteKey = *siteKeyPtr
	adminAPIKey = *adminAPIKeyPtr
	tcUsername = *tcuserPtr
	tcPassword = *tcpasswdPtr

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/payment/register", handleRegisterCC)
	http.HandleFunc("/payment/pay", handlePayment)
	http.HandleFunc("/payment/repeat_pay", handleRepeatPayment)
	http.HandleFunc("/user/info", getUserInformation)
	log.Fatal(http.ListenAndServe(*addrPtr, nil))
}
