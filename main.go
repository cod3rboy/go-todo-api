package main

import (
	"net/http"
	"os"
)

const DefaultPort = "8080"

func main() {
	var port = ""
	if port = os.Getenv("PORT"); port == "" {
		port = DefaultPort
	}
	registerRoutes()
	http.ListenAndServe(":"+port, router)
}
