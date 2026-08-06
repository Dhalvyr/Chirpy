package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(resp http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type validation struct {
		Status string `json:"cleaned_body"`
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

	splittedChirp := strings.Split(params.Body, " ")
	for i := range splittedChirp {
		if _, ok := badWords[strings.ToLower(splittedChirp[i])]; ok {
			splittedChirp[i] = "****"
		}
	}
	finalChirp := strings.Join(splittedChirp, " ")
	validChirp := validation{Status: finalChirp}

	respondWithJSON(resp, 200, validChirp)
	return
}

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
}