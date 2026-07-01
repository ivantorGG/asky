// cmd/migrate/main.go

package main

import (
	"fmt"
	"log"
	"os"

	"asky/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate [up|down|reup|version]")
	}

	cfg := config.Load()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {

	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migrations applied")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migrations rolled back")

	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Version: %d dirty=%v\n", v, dirty)

	default:
		log.Fatal("unknown command [up|down|reup|version]")
	}
}