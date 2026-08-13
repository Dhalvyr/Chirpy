package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/Dhalvyr/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	mux := http.NewServeMux()
	
	apiCfg := apiConfig{db: dbQueries, platform:os.Getenv("PLATFORM")}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	server := &http.Server{
	Addr: ":8080",
	Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
