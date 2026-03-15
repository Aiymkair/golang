package main

import (
	"Assignment-5/db"
	"Assignment-5/internal/config"
	handlers "Assignment-5/internal/handler"
	"Assignment-5/internal/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	repo := repository.New(database)
	h := handlers.New(repo)

	r := mux.NewRouter()
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id1}/common-friends/{id2}", h.GetCommonFriends).Methods("GET")
	r.HandleFunc("/users/{id}", h.SoftDeleteUser).Methods("DELETE")

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
