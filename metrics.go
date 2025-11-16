package main

import (
	"fmt"
	"net/http"
)

const adminPage = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	hitCount := cfg.fileserverHits.Load()

	response := fmt.Sprintf(adminPage, hitCount)

	w.Write([]byte(response))
}
