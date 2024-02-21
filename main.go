package main

import (
	"log"
	"net/http"
)

// Temporary home handler
func home(w http.ResponseWriter, e *http.Request) {
	w.Write([]byte("Hello from snippetbox!"))
}

func main() {
	mux := http.NewServeMux()
	// `/` is a catch-all URL pattern, so all requests will be routed to it
	mux.HandleFunc("/", home)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
