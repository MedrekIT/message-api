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

	mu.Handle("GET /", http.FileServer(http.Dir("./web/static")))
	mu.HandleFunc("GET /api/users", apiCfg.getUsersHandler)
	// mu.HandleFunc("GET /api/users/{userID}", apiCfg.getUserHandler)
	// mu.HandleFunc("GET /api/groups", apiCfg.getGroupsHandler)
	// mu.HandleFunc("GET /api/groups/{groupID}", apiCfg.getGroupHandler)
	// mu.HandleFunc("GET /api/groups/{groupID}/members", apiCfg.getMembersHandler)
	mu.HandleFunc("GET /api/friendships", apiCfg.getFriendsHandler)
	// mu.HandleFunc("GET /api/friendships/requests", apiCfg.getFriendRequestsHandler)

	mu.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mu.HandleFunc("POST /api/users", apiCfg.addUserHandler)
	mu.HandleFunc("POST /api/groups", apiCfg.createGroupHandler)
	// mu.HandleFunc("POST /api/groups/members", apiCfg.joinGroupHandler)
	mu.HandleFunc("POST /api/groups/{groupID}/invitations", apiCfg.createInvitationHandler)
	mu.HandleFunc("POST /api/friendships/{senderLogin}", apiCfg.acceptFriendHandler)
	mu.HandleFunc("POST /api/friendships/requests/{receiverLogin}", apiCfg.requestFriendHandler)
	mu.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mu.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

	return mu
}
