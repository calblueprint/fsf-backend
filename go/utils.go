package main

import (
    "C"
    "fmt"
    "unsafe"
    "reflect"
    "strings"
    "math/rand"
    "net/http"
)

func cArrayToSlice(array unsafe.Pointer, len int) []*C.char {
    var list []*C.char
    sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
    sliceHeader.Cap = len
    sliceHeader.Len = len
    sliceHeader.Data = uintptr(array)
    return list
}

func processTCResponse(resp string) (map[string]string, error) {
    processSingleResp := func(entry string) (string, string, error) {
        items := strings.Split(entry, "=")
        if len(items) == 1 || len(items) > 2 {
            return "", "", fmt.Errorf("malformed entry: %v", entry)
        }

        return items[0], items[1], nil
    }

    entries := strings.Split(resp, "\n")

    ret := make(map[string]string)

    for _, entry := range entries {
        if entry == "" {
            continue
        }

        key, value, err := processSingleResp(entry)
        if err != nil {
            return nil, err
        }

        ret[key] = value
    }

    return ret, nil
}

func getRandomKey() string {
    // get a random key with length 32
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


