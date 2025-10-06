package api

import (
	"net/http"

	"github.com/MedrekIT/message-api/internal/database"
)

type ApiConfig struct {
	Db   *database.Queries
	Port string
}

func Routes(apiCfg *ApiConfig) http.Handler {
	mu := http.NewServeMux()

	mu.Handle("/", http.FileServer(http.Dir("./web/static")))
	//mu.HandleFunc("GET /api/users", apiCfg.getUsersHandler)
	//mu.HandleFunc("GET /api/users/{userID}", apiCfg.getUserHandler)

	mu.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mu.HandleFunc("POST /api/users", apiCfg.addUserHandler)

	return mu
}
