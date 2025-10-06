package main

import (
	"log"
	"fmt"
	"os"
	"net/http"

	"github.com/MedrekIT/message-api/web/api"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	apiCfg := api.ApiConfig{}
	apiCfg.Port = fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

	server := &http.Server{
		Addr: apiCfg.Port,
		Handler: api.Routes(&apiCfg),
	}
	defer server.Close()

	log.Printf("Server running and listening on port %s\n", apiCfg.Port[1:])
	log.Fatal(server.ListenAndServe())
}
