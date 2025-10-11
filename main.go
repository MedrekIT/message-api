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
		log.Fatalf("Error opening database connection - %v\n", err)
	}
	defer func() {
		db.Close()
		log.Println("Database connection closed")
	}()

	apiCfg := api.ApiConfig{
		Port:      fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Db:        database.New(db),
		SecretJWT: secretJWT,
	}

	server := &http.Server{
		Addr:    apiCfg.Port,
		Handler: api.Routes(&apiCfg),
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(time.Hour * 12)
		defer ticker.Stop()
		errorsCounter := 0

		for {
			select {
			case <-ctx.Done():
				log.Println("Cleanup done")
				return
			case <-ticker.C:
				err = automated.DbCleanup(ctx, apiCfg.Db)
				if err != nil && errorsCounter < 3 {
					log.Printf("Error: %v\nRetrying...\n", err)
					errorsCounter++
				} else if err != nil && errorsCounter >= 3 {
					log.Fatalf("Error: %v\n", err)
				} else {
					errorsCounter = 0
				}
			}
		}
	}()

	go func() {
		log.Printf("Server running and listening on port %s\n", apiCfg.Port[1:])
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Error: %v\n", err)
			}
		}
	}()

	<-ch

	log.Println("Shutting down...")
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error while closing server - %v\n", err)
	}
	log.Println("Server closed")
}
