package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pjsmith404/chirpy/internal/auth"
	"github.com/pjsmith404/chirpy/internal/database"
	"net/http"
	"slices"
	"strings"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
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

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't get bearer token",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Failed to validate JWT",
			fmt.Errorf(`Failed to validate JWT: %v`, err),
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
		database.CreateChirpParams{cleanedBody, userID},
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
	authorId := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error
	if authorId != "" {
		parsedAuthorId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Not a valid UUID", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), parsedAuthorId)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}
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

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpId"))
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Not a valid ID",
			err,
		)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusOK, Chirp(chirp))
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpId"))
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Not a valid ID",
			err,
		)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't get bearer token",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Failed to validate JWT",
			fmt.Errorf(`Failed to validate JWT: %v`, err),
		)
		return
	}

	if userID != chirp.UserID {
		respondWithError(
			w,
			http.StatusForbidden,
			"You're not allowed to delete that chirp",
			fmt.Errorf(`JWT user ID %v does not match chirp user ID %v`, userID, chirp.UserID),
		)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to delete chirp",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
