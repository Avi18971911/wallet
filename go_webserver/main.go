package main

import (
	"log"
	"net/http"
)

func main() {
	r := createRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
