// cmd/server/main.go

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

	h := handler.New()

	r := router.New(h)

	log.Printf("Click: http://127.0.0.1:%s\n", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
