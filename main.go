package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server...")
	log.Fatal(s.ListenAndServe())
}
