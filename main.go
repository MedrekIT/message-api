package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MedrekIT/message-api/internal/database"
	"github.com/MedrekIT/message-api/web/api"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("\nError opening database connection - %v\n", err)
	}
	defer db.Close()

	apiCfg := api.ApiConfig{
		Port: fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Db:   database.New(db),
	}

	server := &http.Server{
		Addr:    apiCfg.Port,
		Handler: api.Routes(&apiCfg),
	}
	defer server.Close()

	log.Printf("Server running and listening on port %s\n", apiCfg.Port[1:])
	log.Fatal(server.ListenAndServe())
}
