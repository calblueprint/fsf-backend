package main

import (
    "log"
    "flag"
    "net/http"
    "encoding/json"
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
    apiKey, err := getAPIKey(id)
    if err != nil {
        log.Println(err.Error())
        writeError(w, "Server error when interacting with CiviCRM")
        return
    }

    // write apiKey back in response
    enc := json.NewEncoder(w)

    var key struct {
        Key string `json:"key"`
        Id string `json:"id"`
    }

    key.Key = apiKey
    key.Id = id

    enc.Encode(key)
}

func main() {
    // CAS mobile login server

    addrPtr := flag.String("addr", "0.0.0.0:8080", "address to listen")
    siteKeyPtr := flag.String("sitekey", "", "provide site key of CiviCRM")
    adminApiKeyPtr := flag.String("adminkey", "", "provide the admin API key for CiviCRM")

    flag.Parse()

    siteKey = *siteKeyPtr
    adminApiKey = *adminApiKeyPtr

    http.HandleFunc("/login", handleLogin)
    log.Fatal(http.ListenAndServe(*addrPtr, nil))
}
