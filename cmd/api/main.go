package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"assignment-1/internal/handlers"
	"assignment-1/internal/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handlers.TasksHandler)

	handler := middleware.Logger(middleware.Auth(mux))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server stopped gracefully")
}
