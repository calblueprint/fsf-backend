package main

// #cgo LDFLAGS: ../third_party/tclink/libtclink.a -lssl -lcrypto -lnsl
// #cgo CFLAGS: -I../third_party/tclink/
// #include <stdlib.h>
// #include "tc_wrap.h"
import "C"

import (
    "unsafe"
    "fmt"
)

type TransactionMgr struct {
    CustId string
    Password string
}

func NewTransactionMgr(custId, password string) *TransactionMgr {
    return &TransactionMgr {
        CustId: custId,
        Password: password,
    }
}

// return billing id, err
func (mgr *TransactionMgr) createBillingId(name, ccNumber, expiry, zip string) (string, error) {
    // malloc a C array of char*
    cKeyArray := C.malloc(C.size_t(C.int(7)) * C.size_t(unsafe.Sizeof(uintptr(0))))
    cValueArray := C.malloc(C.size_t(C.int(7)) * C.size_t(unsafe.Sizeof(uintptr(0))))

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
    values[2] = C.CString("store")

    keys[3] = C.CString("name")
    values[3] = C.CString(name)

    keys[4] = C.CString("cc")
    values[4] = C.CString(ccNumber)

    keys[5] = C.CString("exp")
    values[5] = C.CString(expiry)

    keys[6] = C.CString("zip")
    values[6] = C.CString(zip)

    for i := 0; i < 7; i++ {
        defer C.free(unsafe.Pointer(keys[i]))
        defer C.free(unsafe.Pointer(values[i]))
    }

    // allocate buffer for return value
    buf := C.malloc(C.sizeof_char * 1024)
    defer C.free(buf)

    C.TCRequest((**C.char)(cKeyArray), (**C.char)(cValueArray), C.size_t(C.int(7)),
        (*C.char)(buf), C.size_t(C.int(1024)))

    resp := C.GoString((*C.char)(buf))

    respMap, err := processTCResponse(resp)
    if err != nil {
        return "", err
    }

    billingId, ok := respMap["billingid"]
    if !ok {
        return "", fmt.Errorf("billingid is not found in response: %v", resp)
    }

    return billingId, nil
}




