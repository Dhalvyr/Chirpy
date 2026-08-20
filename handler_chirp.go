package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Dhalvyr/Chirpy/internal/auth"
	"github.com/Dhalvyr/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPostChirp(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body     string `json:"body"`
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
	
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, 401, "Unauthorized")
		return
	}

	validatedID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(resp, 401, "Unauthorized")
		return
	}

	newchirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   joinChirp,
		UserID: validatedID,
	})
	if err != nil {
		respondWithError(resp, 500, err.Error())
		return
	}

	finalChirp := Chirp{
		ID:        newchirp.ID,
		CreatedAt: newchirp.CreatedAt,
		UpdatedAt: newchirp.UpdatedAt,
		Body:      newchirp.Body,
		UserID:    newchirp.UserID,
	}

	respondWithJSON(resp, 201, finalChirp)
}

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}
