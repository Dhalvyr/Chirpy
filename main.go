package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	
	apiCfg := apiConfig{}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics", apiCfg.handlerHits)
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	server := &http.Server{
	Addr: ":8080",
	Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(resp http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	msg := fmt.Sprintf("Hits: %d", hits)
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.Write([]byte(msg))
}

func (cfg *apiConfig) handlerReset(resp http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}