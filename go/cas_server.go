package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
)

// global variables to be set by users passing in env variables
var tcUsername, tcPassword string

/*
 * handles a request for a repeated payment, i.e. using billing id
 * requires a POST request with json payload with the following format
 * {
 *  "billingid": "slvkdfjasdoihgjosa",
 *  "amount": "110",
 *	 "email": "some valid email",
 *  "apikey": "some api key"
 * }
 *
 * when successful, returns a json like the following:
 * {
 *   "transid": "a transaction id from TrustCommerce",
 *   "status": "status of transaction { approved, declined, baddata, error }"
 *   "authcode": "auth code for the transaction"
 * }
 */
func handleRepeatPayment(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var ccInfo struct {
		BillingID string `json:"billingid"`
		Amount    string `json:"amount"`
		Email     string `json:"email"`
		ApiKey    string `json:"apikey"`
	}

	if err := dec.Decode(&ccInfo); err != nil {
		writeError(w, "Cannot parse request body correctly")
		return
	}

	// input validation
	var err error
	if ccInfo.BillingID == "" {
		err = errors.New("missing billing id field")
	} else if ccInfo.Amount == "" {
		err = errors.New("missing amount field")
	} else if ccInfo.Email == "" {
		err = errors.New("missing email field")
	} else if ccInfo.ApiKey == "" {
		err = errors.New("missing apiKey field")
	}

	if err != nil {
		log.Println(err.Error())
		writeError(w, err.Error())
		return
	}

	// creates the transaction
	mgr := NewTransactionMgr(tcUsername, tcPassword)
	saleResp, err := mgr.createSaleFromBillingID(ccInfo.BillingID, ccInfo.Amount)
	if err != nil {
		log.Println(err.Error())
		writeError(w, "Payment failed")
		return
	}

	/*
	 * write an error if the transaction is not approved,
	 * else record the contribution in CiviCRM
	 */
	if saleResp.Status != "approved" {
		log.Println(err.Error())
		writeError(w, "transaction not successfully approved")
		return
	} else {
		err := recordTransactionInCiviCRM(ccInfo.Email, ccInfo.ApiKey, saleResp.TransID, ccInfo.Amount)
		if err != nil {
			log.Println(err.Error())
			writeError(w, err.Error())
			return
		}
	}

	// write billing id back in response
	enc := json.NewEncoder(w)

	// see TCSaleResp struct for json response struct
	enc.Encode(saleResp)
}

/*
 * handles a request for a single payment
 * requires a POST request with json payload with the following format
 * {
 *  "name": "John Smith",
 *  "cc": "4111111111111111",
 *  "exp": "0404",
 *  "amount": "110",
 *  "email": "some valid email",
 *  "apikey": "some valid api key"
 * }
 *
 * when successful, returns a json like the following:
 * {
 *   "transid": "a transaction id from TrustCommerce",
 *   "status": "status of transaction { approved, declined, baddata, error }"
 *   "authcode": "auth code for the transaction"
 * }
 */
func handlePayment(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var ccInfo struct {
		Name   string `json:"name"`
		Cc     string `json:"cc"`
		Cvv    string `json:"cvv"`
		Exp    string `json:"exp"`
		Amount string `json:"amount"`
		Email  string `json:"email"`
		ApiKey string `json:"apikey"`
	}

	if err := dec.Decode(&ccInfo); err != nil {
		log.Println(err.Error())
		writeError(w, "Cannot parse request body correctly")
		return
	}

	// input validation
	var err error
	if ccInfo.Name == "" {
		err = errors.New("missing name field")
	} else if ccInfo.Cc == "" {
		err = errors.New("missing cc field")
	} else if ccInfo.Exp == "" {
		err = errors.New("missing exp field")
	} else if ccInfo.Amount == "" {
		err = errors.New("missing amount field")
	} else if ccInfo.Email == "" {
		err = errors.New("missing email field")
	} else if ccInfo.ApiKey == "" {
		err = errors.New("missing apiKey field")
	} else if ccInfo.Cvv == "" {
		err = errors.New("missing cvv field")
	}

	if err != nil {
		log.Println(err.Error())
		writeError(w, err.Error())
		return
	}

	// creates the transaction
	mgr := NewTransactionMgr(tcUsername, tcPassword)

	verifyResp, err := mgr.createVerificationFromCC(ccInfo.Name, ccInfo.Cc, ccInfo.Cvv, ccInfo.Exp, ccInfo.Amount)
	if err != nil {
		log.Println(err.Error())
		writeError(w, "Server side credit-card validation failed")
		return
	}

	saleResp, err := mgr.createSaleFromCC(ccInfo.Name, ccInfo.Cc, ccInfo.Cvv, ccInfo.Exp, ccInfo.Amount)
	// saleResp, err := mgr.createSaleFromCC(ccInfo.Name, ccInfo.Cc, ccInfo.Exp, ccInfo.Amount)
	if err != nil {
		log.Println(err.Error())
		writeError(w, "Payment failed")
		return
	}

	if verifyResp.Status != "approved" {
		log.Println(err.Error())
		writeError(w, "credit card validation transaction not successfully approved")
		return
	}

	/*
	 * var TCSaleResp struct {
	 *		TransID  string `json:"transid"`
	 *		Status   string `json:"status"`
	 *		AuthCode string `json:"authcode"`
	 *	}
	 *
	 *	write an error if the transaction is not approved,
	 *  else record the contribution in CiviCRM
	 */

	if saleResp.Status != "approved" {
		log.Println(err.Error())
		writeError(w, "sale transaction not successfully approved")
		return

	} else {
		/*
				ccInfo struct {
					Name   string `json:"name"`
					Cc     string `json:"cc"`
					Exp    string `json:"exp"`
					Amount string `json:"amount"`
				  Email string `json:"email"`
			  	ApiKey string `json:"apikey"`
				}
		*/
		err := recordTransactionInCiviCRM(ccInfo.Email, ccInfo.ApiKey, saleResp.TransID, ccInfo.Amount)
		/*
			TODO:
			Prevent scenario of payment made, but transaction recorded wrongly; implement either:
				1. repeated retry on failure
				2. separate logs for these scenarios that require admin attention
		*/
		if err != nil {
			log.Println(err.Error())
			writeError(w, err.Error())
			return
		}
	}

	// write billing id back in response
	enc := json.NewEncoder(w)
	// see TCSaleResp struct for json response struct
	enc.Encode(saleResp)
}

/*
	handles a request to store credit card info for repeating payments
	requires a POST request with json payload with the following format
	{
		"name": "John Smith",
		"cc": "4111111111111111",
		"exp": "0404",
		"zip": "90000"
	}

	when success, returns a json like the following:
	{"billingid": "a billing id from TrustCommerce"}
*/

func handleRegisterCC(w http.ResponseWriter, req *http.Request) {
	/*
			requires a POST request with json payload with the following format

		    {
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

	// creates the BillingId
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

/*
	handles a login request
	accepts JSON POST request with the following body:
		{"st": "service token"}

	returns either a HTTP error or a json response like:
		{
			"key": "api key for CiviCRM",
			"id": "contact id for CiviCRM",
			"email": "email address"
		}

	at this time,THE ID TO TALK TO CIVICRM is the users' email address
*/
func handleLogin(w http.ResponseWriter, req *http.Request) {
	/*
	 * requires a POST request with json payload with the following format
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

/*
	handles a userInfo request
	accepts JSON POST request with the following body:
		{
			"key": "api key for CiviCRM",
			"id": "contact id for CiviCRM",
			"email": "email address"
		}
	returns either a HTTP error or a json response like:
		{
			"firstname": "user first name",
			"lastname": "user last name",
			"address": "user address",
			"email": "email address"
		}
*/
func getUserInformation(w http.ResponseWriter, req *http.Request) {
	// requires a POST request

	if req.Method != "POST" {
		writeError(w, "Only POST requests are supported")
		return
	}

	dec := json.NewDecoder(req.Body)
	var key struct {
		Key   string `json:"key"`
		ID    string `json:"id"` // contact_id from CiviCRM
		Email string `json:"email"`
	}

	if err := dec.Decode(&key); err != nil {
		log.Println(dec)
		log.Println(err)
		writeError(w, "Cannot parse request body correctly")
		return
	}

	userInfo, err := getUserInfo(key.Key, key.ID)

	if err != nil {
		log.Println(err.Error())
		writeError(w, "Server error when interacting with CiviCRM")
		return
	}

	// write UserInfo
	enc := json.NewEncoder(w)
	enc.Encode(userInfo)
}

func main() {
	// CAS mobile login server
	addrPtr := flag.String("addr", "0.0.0.0:8000", "address to listen")
	siteKeyPtr := flag.String("sitekey", "", "provide site key of CiviCRM")
	adminAPIKeyPtr := flag.String("adminkey", "", "provide the admin API key for CiviCRM")
	tcuserPtr := flag.String("tcuser", "", "username for Trust Commerce")
	tcpasswdPtr := flag.String("tcpasswd", "", "password for Trust Commerce")

	flag.Parse()

	siteKey = *siteKeyPtr
	adminAPIKey = *adminAPIKeyPtr
	tcUsername = *tcuserPtr
	tcPassword = *tcpasswdPtr

	// endpoints
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/payment/register", handleRegisterCC)
	http.HandleFunc("/payment/pay", handlePayment)
	http.HandleFunc("/payment/repeat_pay", handleRepeatPayment)
	http.HandleFunc("/user/info", getUserInformation)
	log.Fatal(http.ListenAndServe(*addrPtr, nil))
}
