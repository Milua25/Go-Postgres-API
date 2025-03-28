package main

import (
	"fmt"
	"github.com/Golang-Personal-Projects/Go-Projects/06-Go-Postgres-API/router"
	"log"
	"net/http"
)

// using other libraries to
func main() {

	r := router.Router()

	// start server on port 8080
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
