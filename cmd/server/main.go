package main

import (
	"log"
	"net/http"

	"asky/internal/config"
	"asky/internal/database"
	"asky/internal/handler"
	"asky/internal/router"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	h := handler.New(db)

	r := router.New(h)

	log.Printf("Listening on http://127.0.0.1:%s", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}