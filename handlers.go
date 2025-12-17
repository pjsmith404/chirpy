package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool   `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
	}

	respBody := returnVals{
		Valid: true,
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respBody = returnVals{
			Error: "Something went wrong",
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(params.Body) > 140 {
		respBody = returnVals{
			Error: "Chirp is too long",
		}
		w.WriteHeader(http.StatusBadRequest)
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error mashalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}
