package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"slices"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string   `json:"cleaned_body"`
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

	respondWithJson(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})
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
