package main

import (
    "log"
    "flag"
    "net/http"
    //"net/url"
    "encoding/json"
    "encoding/xml"
)

var siteKey string
var adminApiKey string

func writeError(w http.ResponseWriter, msg string) {
    http.Error(w, msg, 500)
}

func writeAccessDenied(w http.ResponseWriter, msg string) {
    http.Error(w, msg, 404)
}

// return
//   bool: whether auth is successful
//   string: id string to retrieve user auth info from civicrm
//   error: if err != nil, there is a server error
func validateToken(token string) (bool, string, error) {
    c := &http.Client{}

    requestUrlPrefix := "https://cas.fsf.org/serviceValidate?service=https://crmserver3d.fsf.org/associate/account&ticket="

    requestUrl := requestUrlPrefix + token

    resp, err := c.Get(requestUrl)

    if (err != nil) {
        return false, "", err
    }

    defer resp.Body.Close()

    if (resp.StatusCode >= 500) {
        return false, "", err
    }

    if (resp.StatusCode >= 400) {
        return false, "", nil
    }

    // resp StatusCode should be 200, parse response for identity

    dec := xml.NewDecoder(resp.Body)

    var r struct {
        XMLName xml.Name `xml:"serviceResponse"`
        AuthSuccess struct {
            XMLName xml.Name `xml:"authenticationSuccess"`
            User string `xml:"user"`
            Attributes struct {
                XMLName xml.Name `xml:"attributes"`
                Email string `xml:"mail"`
            } `xml:"attributes"`
        } `xml:"authenticationSuccess"`
    }

    if err := dec.Decode(&r); err != nil {
        log.Fatal("Failure parsing CAS validate response")
        return false, "", nil
    }

    return true, r.AuthSuccess.Attributes.Email, nil
}

func getAPIKey(id string) (string, error) {
    // TODO: implement secured version
    // current version returns an admin key
    return adminApiKey, nil

    /* WIP */
    /*
    c := &http.Client{}

    requestUrl := "https://crmserver3d.fsf.org/sites/all/modules/civicrm/extern/rest.php"

    var query struct {
        sequential int `json:"sequential"`
        email string `json:"email"`
    }
    query.sequential = 1
    query.email = id

    queryJson, err := json.Marshal(query)
    if err != nil {
        log.Fatal("Error constructing query json for civicrm")
        return "", err
    }

    v := url.Values{}
    v.Add("entity", "Contct")
    v.Add("action", "get")
    v.Add("api_key", adminApiKey)
    v.Add("key", siteKey)
    v.Add("json", string(queryJson))

    resp, err := c.PostForm(requestUrl, v)
    if err != nil {
        return "", err
    }
    */

}


func handleLogin(w http.ResponseWriter, req *http.Request) {
    /** requires a POST request with json payload with the following format
    *
    *  {"st": "service token"}
    */
    if (req.Method != "POST") {
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
