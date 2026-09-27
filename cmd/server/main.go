package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/RoamingAlone/taskforge/internal/config"
	"github.com/RoamingAlone/taskforge/internal/database"
	"github.com/RoamingAlone/taskforge/internal/user"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.New(
		ctx,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Connected to PostgreSQL")

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	router := chi.NewRouter()

	router.Get("/register", userHandler.RegisterForm)
	router.Post("/register", userHandler.RegisterSubmit)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "TaskForge is running!")
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			http.Error(
				w,
				"database unavailable",
				http.StatusServiceUnavailable,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	addr := ":" + cfg.Port

	log.Printf("TaskForge running on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
