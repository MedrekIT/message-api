package api

import (
	"net/http"

	"github.com/MedrekIT/message-api/internal/database"
)

type ApiConfig struct {
	Db        *database.Queries
	Port      string
	SecretJWT string
}

func Routes(apiCfg *ApiConfig) http.Handler {
	mu := http.NewServeMux()

	mu.Handle("/", http.FileServer(http.Dir("./web/static")))
	//mu.HandleFunc("GET /api/users", apiCfg.getUsersHandler)
	//mu.HandleFunc("GET /api/users/{userID}", apiCfg.getUserHandler)

	mu.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mu.HandleFunc("POST /api/users", apiCfg.addUserHandler)
	mu.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mu.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

	return mu
}
