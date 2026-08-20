package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dhalvyr/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Duration int `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	standardDuration := 3600

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}

	if params.Duration > 0 && params.Duration < 3600 {
		standardDuration = params.Duration
	}

	fetchUser, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(resp, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, fetchUser.HashedPassword)
	if err != nil || !match {
		respondWithError(resp, 401, "Incorrect email or password")
		return
	}

	newToken, err := auth.MakeJWT(fetchUser.ID, cfg.jwtSecret, time.Duration(standardDuration) * time.Second)
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}
	returnUser := response{
		User: User{
			ID:        fetchUser.ID,
			CreatedAt: fetchUser.CreatedAt,
			UpdatedAt: fetchUser.UpdatedAt,
			Email:     fetchUser.Email,
		},
		Token: newToken,
	}

	respondWithJSON(resp, 200, returnUser)
}