package main

import (
	"encoding/json"
	"net/http"

	"github.com/Dhalvyr/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, 400, err.Error())
		return
	}

	fetchUser, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(resp, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, fetchUser.HashedPassword)
	if err != nil || match == false {
		respondWithError(resp, 401, "Incorrect email or password")
		return
	}

	returnUser := User{
		ID:        fetchUser.ID,
		CreatedAt: fetchUser.CreatedAt,
		UpdatedAt: fetchUser.UpdatedAt,
		Email:     fetchUser.Email,
	}

	respondWithJSON(resp, 200, returnUser)

}