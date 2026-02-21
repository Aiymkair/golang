package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assignment-2/internal/delivery/handler"
	"assignment-2/internal/repository"
	"assignment-2/internal/repository/_postgres"
	"assignment-2/internal/usecase"
	"assignment-2/pkg/modules"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	dbConfig := &modules.PostgreSQL{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "1234",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}

	ctx := context.Background()
	pgDialect := _postgres.NewPGXDialect(ctx, dbConfig)
	_postgres.AutoMigrate(dbConfig)

	repos := repository.NewRepositories(pgDialect)

	userUsecase := usecase.NewUserUsecase(repos.Users)

	// Создаём обработчик – переменная h, не перекрывает пакет handler
	h := handler.NewHandler(userUsecase)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// Вызываем функции из пакета handler (не из переменной h)
	r.Use(handler.LoggingMiddleware)
	r.Use(handler.APIKeyMiddleware)

	r.Get("/health", h.Healthcheck)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.GetAllUsers)
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUserByID)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUser)
	})

	// 7. Запуск сервера с graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server stopped")
}
