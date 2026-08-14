package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetSingleChirp(resp http.ResponseWriter, req *http.Request) {
	reqChirpID := req.PathValue("chirpID")
	parsedChirp, err := uuid.Parse(reqChirpID)
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}

	reqChirp, err := cfg.db.GetSingleChirp(req.Context(), parsedChirp)
	if err != nil {
		respondWithError(resp, 404, err.Error())
		return
	}

	finalChirp := Chirp{
		ID:        reqChirp.ID,
		CreatedAt: reqChirp.CreatedAt,
		UpdatedAt: reqChirp.UpdatedAt,
		Body:      reqChirp.Body,
		UserID:    reqChirp.UserID,
	}

	respondWithJSON(resp, 200, finalChirp)
}
