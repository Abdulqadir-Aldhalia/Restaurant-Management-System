package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"server-side/controllers"
	"server-side/services"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Using the current working directory instead of executable path
	migrationPath := "file://" + services.GetRootpath("database/migrations")
	fmt.Println("Migration path:", migrationPath)

	mig, err := migrate.New(
		migrationPath,
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := mig.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
		log.Printf("migrations: %s", err.Error())
	}

	addr := flag.String("addr", ":8000", "HTTP Networking address")
	flag.Parse()

	services.SetDB(db)
	controllers.SetServerPort(*addr)

	controllers.Controllers()
}
