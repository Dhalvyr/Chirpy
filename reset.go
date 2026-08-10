package main

import "net/http"

func (cfg *apiConfig) handlerReset(resp http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(resp, 403, "Forbidden")
		return
	}
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}
	cfg.fileserverHits.Store(0)
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}