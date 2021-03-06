package main

// #cgo LDFLAGS: ../third_party/tclink/libtclink.a -lssl -lcrypto -L/usr/local/opt/openssl/lib
// #cgo CFLAGS: -I../third_party/tclink/
// #include <stdlib.h>
// #include "tc_wrap.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type TransactionMgr struct {
	CustId   string
	Password string
}

type TCSaleResp struct {
	TransID  string `json:"transid"`
	Status   string `json:"status"`
	AuthCode string `json:"authcode"`
}

type TCVerifyResp struct {
	TransID  string `json:"transid"`
	Status   string `json:"status"`
	AuthCode string `json:"authcode"`
	Avs      string `json:"avs"`
}

var avsResponseCodes = map[string]string{
	"A": "Street address matches, but five-digit and nine-digit postal code do not match.",
	"B": "Street address matches, but postal code not verified.",
	"C": "Street address and postal code do not match.",
	"D": "Street address and postal code match.",
	"E": "AVS data is invalid or AVS is not allowed for this card type.",
	"F": "Card member's name does not match, but billing postal code matches.",
	"G": "Non-U.S. issuing bank does not support AVS.",
	"H": "Card member's name does not match. Street address and postal code match.",
	"I": "Address not verified.",
	"J": "Card member's name, billing address, and postal code match.",
	"K": "Card member's name matches but billing address and billing postal code do not match.",
	"L": "Card member's name and billing postal code match, but billing address does not match.",
	"M": "Street address and postal code match.",
	"N": "Street address and postal code do not match.",
	"O": "Card member's name and billing address match, but billing postal code does not match.",
	"P": "Postal code matches, but street address not verified.",
	"Q": "Card member's name, billing address, and postal code match.",
	"R": "System unavailable.",
	"S": "Bank does not support AVS.",
	"T": "Card member's name does not match, but street address matches.",
	"U": "Address information unavailable.",
	"V": "Card member's name, billing address, and billing postal code match.",
	"W": "Street address does not match, but nine-digit postal code matches.",
	"X": "Street address and nine-digit postal code match.",
	"Y": "Street address and five-digit postal code match.",
	"Z": "Street address does not match, but five-digit postal code matches.",
}

func NewTransactionMgr(custId, password string) *TransactionMgr {
	return &TransactionMgr{
		CustId:   custId,
		Password: password,
	}
}

// return transaction status struct, err
func (mgr *TransactionMgr) createSaleFromCC(name, ccNumber, cvv, expiry, amount string) (*TCSaleResp, error) {
	// malloc a C array of char*
	mapSize := 9

	cKeyArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cValueArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	defer C.free(cKeyArray)
	defer C.free(cValueArray)

	// convert C array to go slice for addressing
	keys := cArrayToSlice(cKeyArray, mapSize)
	values := cArrayToSlice(cValueArray, mapSize)

	// set parameters
	keys[0] = C.CString("custid")
	values[0] = C.CString(mgr.CustId)

	keys[1] = C.CString("password")
	values[1] = C.CString(mgr.Password)

	keys[2] = C.CString("action")
	values[2] = C.CString("sale")

	keys[3] = C.CString("name")
	values[3] = C.CString(name)

	keys[4] = C.CString("cc")
	values[4] = C.CString(ccNumber)

	keys[5] = C.CString("exp")
	values[5] = C.CString(expiry)

	keys[6] = C.CString("amount")
	values[6] = C.CString(amount)

	keys[7] = C.CString("cvv")
	values[7] = C.CString(cvv)

	keys[8] = C.CString("checkcvv")
	values[8] = C.CString("y")

	for i := 0; i < mapSize; i++ {
		defer C.free(unsafe.Pointer(keys[i]))
		defer C.free(unsafe.Pointer(values[i]))
	}

	// allocate buffer for return value
	buf := C.malloc(C.sizeof_char * 1024)
	defer C.free(buf)

	return mgr.createSaleHelper((**C.char)(cKeyArray), (**C.char)(cValueArray), mapSize,
		(*C.char)(buf), 1024)
}

// return transaction status struct, err
func (mgr *TransactionMgr) createSaleFromBillingID(billingID, amount string) (*TCSaleResp, error) {
	// malloc a C array of char*
	mapSize := 5

	cKeyArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cValueArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	defer C.free(cKeyArray)
	defer C.free(cValueArray)

	// convert C array to go slice for addressing
	keys := cArrayToSlice(cKeyArray, 7)
	values := cArrayToSlice(cValueArray, 7)

	// set parameters
	keys[0] = C.CString("custid")
	values[0] = C.CString(mgr.CustId)

	keys[1] = C.CString("password")
	values[1] = C.CString(mgr.Password)

	keys[2] = C.CString("action")
	values[2] = C.CString("sale")

	keys[3] = C.CString("billingid")
	values[3] = C.CString(billingID)

	keys[4] = C.CString("amount")
	values[4] = C.CString(amount)

	for i := 0; i < mapSize; i++ {
		defer C.free(unsafe.Pointer(keys[i]))
		defer C.free(unsafe.Pointer(values[i]))
	}

	// allocate buffer for return value
	buf := C.malloc(C.sizeof_char * 1024)
	defer C.free(buf)

	return mgr.createSaleHelper((**C.char)(cKeyArray), (**C.char)(cValueArray), mapSize,
		(*C.char)(buf), 1024)
}

// return raw response string, parsed map, and err
func (mgr *TransactionMgr) TCTransactionHelper(cKeyArray **C.char, cValueArray **C.char, mapSize int, dest *C.char, bufSize int) (string, map[string]string, error) {

	C.TCRequest(cKeyArray, cValueArray, C.size_t(C.int(mapSize)), dest, C.size_t(C.int(bufSize)))

	resp := C.GoString((*C.char)(dest))

	respMap, err := processTCResponse(resp)
	if err != nil {
		return "", nil, err
	}

	return resp, respMap, nil
}

// return transaction status struct, err
func (mgr *TransactionMgr) createSaleHelper(cKeyArray **C.char, cValueArray **C.char, mapSize int, dest *C.char, bufSize int) (*TCSaleResp, error) {
	resp, respMap, err := mgr.TCTransactionHelper(cKeyArray, cValueArray, mapSize, dest, bufSize)

	if err != nil {
		return nil, err
	}

	var tcResp TCSaleResp
	var ok bool

	tcResp.TransID, ok = respMap["transid"]
	if !ok {
		return nil, fmt.Errorf("transid is not found in response: %v", resp)
	}
	tcResp.Status, ok = respMap["status"]
	if !ok {
		return nil, fmt.Errorf("status is not found in response: %v", resp)
	}
	tcResp.AuthCode, ok = respMap["authcode"]
	if !ok {
		return nil, fmt.Errorf("authcode is not found in response: %v", resp)
	}

	return &tcResp, nil
}

// return billing id, err
func (mgr *TransactionMgr) createBillingId(name, ccNumber, expiry, zip string) (string, error) {
	// malloc a C array of char*
	mapSize := 7

	cKeyArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cValueArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	defer C.free(cKeyArray)
	defer C.free(cValueArray)

	// convert C array to go slice for addressing
	keys := cArrayToSlice(cKeyArray, mapSize)
	values := cArrayToSlice(cValueArray, mapSize)

	// set parameters
	keys[0] = C.CString("custid")
	values[0] = C.CString(mgr.CustId)

	keys[1] = C.CString("password")
	values[1] = C.CString(mgr.Password)

	keys[2] = C.CString("action")
	values[2] = C.CString("store")

	keys[3] = C.CString("name")
	values[3] = C.CString(name)

	keys[4] = C.CString("cc")
	values[4] = C.CString(ccNumber)

	keys[5] = C.CString("exp")
	values[5] = C.CString(expiry)

	keys[6] = C.CString("zip")
	values[6] = C.CString(zip)

	for i := 0; i < mapSize; i++ {
		defer C.free(unsafe.Pointer(keys[i]))
		defer C.free(unsafe.Pointer(values[i]))
	}

	// allocate buffer for return value
	buf := C.malloc(C.sizeof_char * 1024)
	defer C.free(buf)

	resp, respMap, err := mgr.TCTransactionHelper((**C.char)(cKeyArray), (**C.char)(cValueArray), mapSize, (*C.char)(buf), 1024)

	if err != nil {
		return "", err
	}

	billingId, ok := respMap["billingid"]
	if !ok {
		return "", fmt.Errorf("billingid is not found in response: %v", resp)
	}

	return billingId, nil
}

// INCOMPLETE VERIFICATION FUNCTIONS

// func (mgr *TransactionMgr) verifyAddress(action, ccNumber, expiry, amount string) (*TCSaleResp, error) {

// }

// return transaction status struct, err
func (mgr *TransactionMgr) createVerificationFromCC(name, ccNumber, expiry, cvv string) (*TCVerifyResp, error) {
	// malloc a C array of char*
	mapSize := 5 //DIFF 1: mapSize

	cKeyArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cValueArray := C.malloc(C.size_t(C.int(mapSize)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	defer C.free(cKeyArray)
	defer C.free(cValueArray)

	// convert C array to go slice for addressing
	keys := cArrayToSlice(cKeyArray, mapSize)
	values := cArrayToSlice(cValueArray, mapSize)

	// set parameters
	keys[0] = C.CString("custid")
	values[0] = C.CString(mgr.CustId)

	keys[1] = C.CString("password")
	values[1] = C.CString(mgr.Password)

	keys[2] = C.CString("action")
	values[2] = C.CString("verify") // DIFF 2: action is verify

	// keys[3] = C.CString("name")
	// values[3] = C.CString(name)

	keys[3] = C.CString("cc")
	values[3] = C.CString(ccNumber)

	keys[4] = C.CString("exp")
	values[4] = C.CString(expiry)

	// keys[5] = C.CString("cvv")
	// values[5] = C.CString(cvv)

	// keys[6] = C.CString("verify")
	// values[6] = C.CString("y")

	for i := 0; i < mapSize; i++ {
		defer C.free(unsafe.Pointer(keys[i]))
		defer C.free(unsafe.Pointer(values[i]))
	}

	// allocate buffer for return value
	buf := C.malloc(C.sizeof_char * 1024)
	defer C.free(buf)

	return mgr.createVerifyCardHelper((**C.char)(cKeyArray), (**C.char)(cValueArray), mapSize,
		(*C.char)(buf), 1024)
}

func (mgr *TransactionMgr) createVerifyCardHelper(cKeyArray **C.char, cValueArray **C.char, mapSize int, dest *C.char, bufSize int) (*TCVerifyResp, error) {
	resp, respMap, err := mgr.TCTransactionHelper(cKeyArray, cValueArray, mapSize, dest, bufSize)

	if err != nil {
		return nil, err
	}

	var tcResp TCVerifyResp
	var ok bool

	tcResp.TransID, ok = respMap["transid"]
	if !ok {
		return nil, fmt.Errorf("transid is not found in response: %v", resp)
	}
	tcResp.Status, ok = respMap["status"]
	if !ok {
		return nil, fmt.Errorf("status is not found in response: %v", resp)
	}
	tcResp.AuthCode, ok = respMap["authcode"]
	if !ok {
		return nil, fmt.Errorf("authcode is not found in response: %v", resp)
	}
	tcResp.Avs, ok = respMap["avs"]
	if !ok {
		return nil, fmt.Errorf("avs is not found in response: %v", resp)
	}
	return &tcResp, nil
}
