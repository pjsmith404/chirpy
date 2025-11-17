package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	handleFs := http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handleFs))
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server...")
	log.Fatal(s.ListenAndServe())
}
