package main

import (
	"net/http"
    "sync/atomic"
    "fmt"
)

type apiConfig struct {
    fileserverHits  atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    res := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
    `
    fmt.Fprintf(w, res, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handleResetMetrics(_ http.ResponseWriter, _ *http.Request) {
    cfg.fileserverHits.Store(0)
}

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

    // Fileserver Handlers
    apiCfg := apiConfig{}
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
    mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets/")))))

    // APIs
    mux.HandleFunc("GET /api/healthz", handleHealthz)
    mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
    mux.HandleFunc("POST /admin/reset", apiCfg.handleResetMetrics)
    
    // serve and listen to connections
    server.ListenAndServe()
}
