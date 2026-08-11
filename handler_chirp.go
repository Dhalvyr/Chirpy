package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/Dhalvyr/Chirpy/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPostChirp(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		PosterID string `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(resp, 400, "Chirp is too long")
		return
	}

	splittedChirp := strings.Split(params.Body, " ")
	for i := range splittedChirp {
		if _, ok := badWords[strings.ToLower(splittedChirp[i])]; ok {
			splittedChirp[i] = "****"
		}
	}
	joinChirp := strings.Join(splittedChirp, " ")
	parsedID, err := uuid.Parse(params.PosterID)
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}

	newchirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: joinChirp,
		UserID: parsedID,
	})
	if err != nil {
		respondWithError(resp, 500, err.Error())
		return
	}

	finalChirp := Chirp{
		ID: newchirp.ID,
		CreatedAt: newchirp.CreatedAt,
		UpdatedAt: newchirp.UpdatedAt,
		Body: newchirp.Body,
		UserID: newchirp.UserID,
	}

	respondWithJSON(resp, 201, finalChirp)
}

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
}