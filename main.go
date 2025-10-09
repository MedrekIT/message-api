package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MedrekIT/message-api/internal/automated"
	"github.com/MedrekIT/message-api/internal/database"
	"github.com/MedrekIT/message-api/web/api"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	secretJWT := os.Getenv("SECRET_JWT")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("\nError opening database connection - %v\n", err)
	}
	defer db.Close()

	apiCfg := api.ApiConfig{
		Port:      fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Db:        database.New(db),
		SecretJWT: secretJWT,
	}

	server := &http.Server{
		Addr:    apiCfg.Port,
		Handler: api.Routes(&apiCfg),
	}
	defer server.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	ticker := time.NewTicker(time.Hour * 12)
	errorsCounter := 0
	go func() {
		for {
			select {
			case <-ch:
				log.Println("Shutting down...")
				cancel()
			case <-ctx.Done():
				log.Println("Cleanup done")
				log.Println("Server closed")
				os.Exit(0)
			case <-ticker.C:
				err = automated.DbCleanup(ctx, apiCfg.Db)
				if err != nil && errorsCounter < 3 {
					log.Printf("\nError: %v\nRetrying...\n", err)
					errorsCounter++
				} else if err != nil && errorsCounter >= 3 {
					log.Fatalf("\nError: %v\n", err)
				}
			}
		}
	}()
	log.Printf("Server running and listening on port %s\n", apiCfg.Port[1:])
	log.Fatal(server.ListenAndServe())
}
