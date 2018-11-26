package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

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

	result, id, err := validateToken(serviceToken)

	if err != nil {
		writeError(w, "Server error when validating token")
		return
	} else if !result {
		writeAccessDenied(w, "Ticket authentication failed")
		return
	}

	// authentication success
	APIKey, err := getAPIKey(id)
	if err != nil {
		log.Println(err.Error())
		writeError(w, "Server error when interacting with CiviCRM")
		return
	}

	// write APIKey back in response
	enc := json.NewEncoder(w)

	var key struct {
		Key string `json:"key"`
		ID  string `json:"id"`
	}

	key.Key = APIKey
	key.ID = id

	enc.Encode(key)
}

func main() {
	// CAS mobile login server

	addrPtr := flag.String("addr", "0.0.0.0:8080", "address to listen")
	siteKeyPtr := flag.String("sitekey", "", "provide site key of CiviCRM")
	adminAPIKeyPtr := flag.String("adminkey", "", "provide the admin API key for CiviCRM")

	flag.Parse()

	siteKey = *siteKeyPtr
	adminAPIKey = *adminAPIKeyPtr

	http.HandleFunc("/login", handleLogin)
	log.Fatal(http.ListenAndServe(*addrPtr, nil))
}
