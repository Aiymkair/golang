package main

import (
	"Assignment-7/internal/app"
	"log"
)

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
