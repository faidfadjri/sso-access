package main

import (
	"log"

	database "akastra-access/internal/infrastructure/databases"
	"akastra-access/internal/infrastructure/databases/seeding"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	seeding.Run(db)
}
