package main

import (
    "log"
    "flag"
    "net/http"
    "net/url"
    "encoding/json"
    "encoding/xml"
    "math/rand"
    "fmt"
)

var siteKey string
var adminApiKey string

func getRandomKey() string {
    bytes := make([]byte, 32)
    for i := 0; i < 32; i++ {
        bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
    }
    return string(bytes)
}

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
    /** Because of the design of CiviCRM, API Key is only shown when we do an update
        i.e. a create with contact_id specified.

        We first query for the contact id that matches this identity, then run an
        update to retrieve the API key.

        If current API key is not set, we set the API key and return to the user.
    */
    c := &http.Client{}

    requestUrl := "https://crmserver3d.fsf.org/sites/all/modules/civicrm/extern/rest.php"

    // Query for contact id
    var idQuery struct {
        Sequential int `json:"sequential"`
        Return string `json:"return"`
        Email string `json:"email"`
    }
    idQuery.Sequential = 1
    idQuery.Email = id
    idQuery.Return = "id"

    idQueryJson, err := json.Marshal(idQuery)
    if err != nil {
        log.Fatal("Error constructing query json for civicrm")
        return "", err
    }

    v := url.Values{}
    v.Add("entity", "Contact")
    v.Add("action", "get")
    v.Add("api_key", adminApiKey)
    v.Add("key", siteKey)
    v.Add("json", string(idQueryJson))

    resp, err := c.PostForm(requestUrl, v)
    if err != nil {
        return "", err
    }
    if resp.StatusCode >= 400 {
        return "", fmt.Errorf("Bad status code: %v", string(resp.StatusCode))
    }

    var idQueryResp struct {
        Error int `json:"is_error"`
        Values []struct {
            Id string `json:"id"`
        } `json:"values"`
    }

    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&idQueryResp)
    if err != nil || len(idQueryResp.Values) != 1 || idQueryResp.Error != 0 {
        return "", fmt.Errorf("Bad response")
    }

    resp.Body.Close()

    contactId := idQueryResp.Values[0].Id
    log.Println("Contact id is:" + contactId)

    // Use an update query to check for API key
    updateQueryJson := `{"id":"` + contactId + `"}`

    v.Set("action", "create")
    v.Set("json", updateQueryJson)

    resp, err = c.PostForm(requestUrl, v)
    if err != nil {
        return "", err
    }
    if resp.StatusCode >= 400 {
        return "", fmt.Errorf("Bad status code: %v", string(resp.StatusCode))
    }

    var updateQueryResp struct {
        Error int `json:"is_error"`
        Values map[string] struct {
            ApiKey string `json:"api_key"`
        } `json:"values"`
    }

    dec = json.NewDecoder(resp.Body)
    err = dec.Decode(&updateQueryResp)
    if err != nil || updateQueryResp.Error != 0 {
        return "", fmt.Errorf("Bad response")
    }

    resp.Body.Close()

    if updateQueryResp.Values[contactId].ApiKey != "" {
        return updateQueryResp.Values[contactId].ApiKey, nil
    }

    log.Println("Found API key:" + updateQueryResp.Values[contactId].ApiKey)
    // api key is not set, need to update API key and return
    newApiKey := getRandomKey()

    updateQueryJson = `{"id":"` + idQueryResp.Values[0].Id +
        `", "api_key": "` + newApiKey + `"}`
    v.Set("json", updateQueryJson)
    resp, err = c.PostForm(requestUrl, v)
    if err != nil {
        return "", err
    }
    if resp.StatusCode >= 400 {
        return "", fmt.Errorf("Bad status code: %v", string(resp.StatusCode))
    }

    dec = json.NewDecoder(resp.Body)
    err = dec.Decode(&updateQueryResp)
    if err != nil || updateQueryResp.Error != 0 {
        return "", fmt.Errorf("Bad response")
    }
    log.Println("Set API key:" + updateQueryResp.Values[contactId].ApiKey)

    return newApiKey, nil
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
