package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"github.com/pjsmith404/chirpy/internal/database"
	"time"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type Chirps []Chirp

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := cleanBody(params.Body)

	chirp, err := cfg.db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{cleanedBody, params.UserId},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create chirp",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusCreated, Chirp(chirp))
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't get chirps",
			err,
		)
		return
	}

	jsonChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		jsonChirps[i] = Chirp(chirp)
	}

	respondWithJson(w, http.StatusOK, jsonChirps)
}

func cleanBody(body string) string {
	splitBody := strings.Split(body, " ")
	cleanedBody := make([]string, 0, len(splitBody))

	bannedWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	for _, word := range splitBody {
		if slices.Contains(bannedWords, strings.ToLower(word)) {
			cleanedBody = append(cleanedBody, "****")
		} else {
			cleanedBody = append(cleanedBody, word)
		}
	}

	return strings.Join(cleanedBody, " ")
}
