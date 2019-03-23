## This folder contains the Go files for the FSF mobile backend.
It houses the code to interact with CiviCRM and TrustCommerce: accepting payments and handling CAS Login. 

### About
`cas_server.go` is our server that is out intemediary between the client app (found in fsf-mobile) and TrustCommerce/CiviCRM. It calls functions from `civicrm_client.go`, which contains functions that query CiviCRM, `login.go`, which contains functions called for the login flow, and `payment.go`, which contains functions called for the payment flow.

### _test.go
`cas_server_test.go` and `civicrm_client_test.go` contain tests for cas_server and login respectively. There are number of tests in them right now and we plan to add more tests in the future, time permitting.

### cas_server.go
`cas_server.go` defines 5 endpoints:
1. /login
2. /payment/register
3. /payment/pay
4. /payment/repeat_pay
5. /user/info

### civicrm_client.go
`civicrm_client.go` contains 5 helper functions:

1. queryCiviCRM queries CiviCRM with an encoded REST query, and stores the information received in a decoded json object
2. getAPIKey gets the API key for a user using a unique ID, which is currently the user's email
3. validateAPIKeyForUpdateRequests validates the user's APIKey which we get from the client with the user's API key from CiviCRM
4. getUserInfo gets the user's info
5. recordTransactionInCiviCRM records a contribution a user made into CiviCRM as a Contribution object

### login.go
validateToken is a function that accepts a service token from the client, who obtains it from
CAS. It interacts with CAS server to validate the token.

### payment.go
`payment.go` contains 6 helper functions:

1. NewTransactionMgr creates and returns a new TransactionMgr
2. createSaleFromCC creates a sale from a credit card number
3. createSaleFromBillingID creates a sale from a billing id
4. TCTransactionHelper returns a raw response string, parsed map, and err
5. createSaleHelper creates and returns a transaction status struct
6. createBillingId creates a billing Id from a transaction
