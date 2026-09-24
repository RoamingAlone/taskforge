package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RoamingAlone/taskforge/internal/config"
	"github.com/go-chi/chi/v5"
)

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

	cfg := config.Load()

	addr := ":" + cfg.Port

	log.Printf("TaskForge running on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}