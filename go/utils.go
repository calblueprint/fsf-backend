package main

import (
    "math/rand"
    "net/http"
)

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


