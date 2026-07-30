package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerReadiness)
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
