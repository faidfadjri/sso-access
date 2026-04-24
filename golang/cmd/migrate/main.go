package main

import (
	"log"

	database "akastra-access/internal/infrastructure/databases"
	"akastra-access/internal/infrastructure/databases/migrations"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Migration terminated with error: %v", err)
	}
}
