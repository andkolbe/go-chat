package main

// A package is a way to group functions, and it's made up of all the files in the same directory
// A main function executes by default when you run the main package

import (
	"log"
	"net/http"

	"github.com/andkolbe/go-chat/internal/handlers"
)


func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}

