package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pjsmith404/chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	db             *database.Queries
	fileserverHits atomic.Int32
	jwtSecret      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db:        dbQueries,
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()
	handleFs := http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handleFs))
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server...")
	log.Fatal(s.ListenAndServe())
}
