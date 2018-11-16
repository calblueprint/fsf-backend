package main

import (
    "log"
    "flag"
    "net/http"
    "encoding/json"
)

var tcUsername, tcPassword string

// handles a request that registers credit card info
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
        Cc string `json:"cc"`
        Exp string `json:"exp"`
        Zip string `json:"zip"`
    }

    if err := dec.Decode(&ccInfo); err != nil {
        writeError(w, "Cannot parse request body correctly")
        return
    }

    mgr := NewTransactionMgr(tcUsername, tcPassword)
    billingId, err := mgr.createBillingId(ccInfo.Name, ccInfo.Cc, ccInfo.Exp, ccInfo.Zip)

    if err != nil {
        log.Println(err.Error())
        writeError(w, "Payment failed")
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
    tcuserPtr := flag.String("tcuser", "", "username for Trust Commerce")
    tcpasswdPtr := flag.String("tcpasswd", "", "password for Trust Commerce")

    flag.Parse()

    siteKey = *siteKeyPtr
    adminApiKey = *adminApiKeyPtr
    tcUsername = *tcuserPtr
    tcPassword = *tcpasswdPtr

    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/payment/register", handleRegisterCC)
    log.Fatal(http.ListenAndServe(*addrPtr, nil))
}
