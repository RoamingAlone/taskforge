package main

import {
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
}

func main() {
	router := chi.NewRouter()
	
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "TaskForge is running!")
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Println("TaskForge running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}