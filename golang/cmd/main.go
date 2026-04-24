package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"akastra-access/internal/app/bootstrap"
	"akastra-access/internal/interface/routes"
)

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	time.Local = loc
}

func main() {
	if os.Getenv("DEBUG") == "true" {
		logFile, err := os.OpenFile(
			"/var/log/akastra-access.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0664,
		)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	deps := bootstrap.InitDependencies()

	router := routes.InitRouter(deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
