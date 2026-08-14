package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(resp http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(resp, 500, err.Error())
		return
	}

	chirpList := []Chirp{}

	for _, chirp := range dbChirps {
		finalChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpList = append(chirpList, finalChirp)
	}

	respondWithJSON(resp, 200, chirpList)
}
