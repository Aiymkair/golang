package main

import (
	"context"
	_ "database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"assignment-4/internal/delivery/handler"
	"assignment-4/internal/repository"
	"assignment-4/internal/repository/_postgres"
	"assignment-4/internal/usecase"
	"assignment-4/pkg/modules"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// getEnv читает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func main() {
	// Конфигурация из переменных окружения
	dbConfig := &modules.PostgreSQL{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		Username:    getEnv("DB_USER", "postgres"),
		Password:    getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "mydb"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: time.Duration(getEnvAsInt("DB_TIMEOUT", 5)) * time.Second,
	}

	// Подключение к БД
	ctx := context.Background()
	pgDialect := _postgres.NewPGXDialect(ctx, dbConfig)

	// Применяем миграции
	_postgres.AutoMigrate(dbConfig)

	// Репозитории
	repos := repository.NewRepositories(pgDialect)

	// Usecase
	userUsecase := usecase.NewUserUsecase(repos.Users)

	// Обработчик HTTP
	h := handler.NewHandler(userUsecase)

	// Роутер
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
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

	// HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Запуск в горутине
	go func() {
		log.Println("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			log.Fatalf("listen error: %s", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gracefully...")

	// Даём серверу 10 секунд на завершение запросов
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Закрываем соединение с БД
	if err := pgDialect.DB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		log.Println("Database connection closed")
	}

	log.Println("Server stopped")
}
