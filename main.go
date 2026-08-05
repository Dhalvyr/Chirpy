package main

import (
	"encoding/json"
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
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
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

func (cfg *apiConfig) handlerMetrics(resp http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.Write([]byte(fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
		hits,
	)))
}

func (cfg *apiConfig) handlerReset(resp http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}

func handlerValidateChirp(resp http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type validation struct {
		Status bool `json:"valid"`
	}

	validChirp := validation{
		Status: true,
	}

	decoder := json.NewDecoder(req.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, 500, err.Error())
		return
	}
	if len(params.Body) > 140 {
		respondWithError(resp, 400, "Chirp is too long")
		return
	}

	respondWithJSON(resp, 200, validChirp)
	return
}

func respondWithError(resp http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	errorText := errorResponse{
		Error: msg,
	}
	respondWithJSON(resp, code, errorText)
	return
}

func respondWithJSON(resp http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		resp.WriteHeader(500)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	resp.Write(data)
	return
}