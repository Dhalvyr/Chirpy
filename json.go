package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
