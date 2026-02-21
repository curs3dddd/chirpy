package main

import (
	"net/http"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

func main() {
    // initialize the server stats
    mux := http.NewServeMux()
    server := &http.Server {
        Addr:       ":8080",
        Handler:    mux,
    }

    // handlers
    mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
    mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets"))))

    // custom handlers
    mux.HandleFunc("/healthz", handleHealthz)
    
    // serve and listen to connections
    server.ListenAndServe()
}
