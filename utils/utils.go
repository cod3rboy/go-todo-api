package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func GenerateRandomID() string {
	const length = 16
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func MethodNotAllowed(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write([]byte("method not allowed"))
}

func BadRequest(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte("bad request"))
}
