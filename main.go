package main

import (
	"net/http"
)

func main() {
    // initialize the server stats
    mux := http.NewServeMux()
    server := &http.Server {
        Addr:       ":8080",
        Handler:    mux,
    }

    // handlers
    mux.Handle("/", http.FileServer(http.Dir(".")))
    mux.Handle("/assets", http.FileServer(http.Dir("./assets/logo.png")))
    
    // serve and listen to connections
    server.ListenAndServe()
}
