package api

import (
	"net/http"
	"time"

	"github.com/MedrekIT/message-api/internal/auth"
	"github.com/MedrekIT/message-api/internal/database"
)

func (apiCfg *ApiConfig) getFriendsHandler(w http.ResponseWriter, r *http.Request) {
}

func (apiCfg *ApiConfig) acceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid ID", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, apiCfg.SecretJWT)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid ID", err)
		return
	}

	user, err := apiCfg.Db.GetUserByID(r.Context(), userID)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid ID", err)
		return
	}

	senderLogin := r.PathValue("senderLogin")
	sender, err := apiCfg.Db.GetUserByLogin(r.Context(), senderLogin)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "User not found", err)
		return
	}

	newAcceptParams := database.AcceptFriendshipParams{
		UserID:     sender.ID,
		ReceiverID: user.ID,
	}
	err = apiCfg.Db.AcceptFriendship(r.Context(), newAcceptParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	successResponse(w, http.StatusAccepted, nil)
}

func (apiCfg *ApiConfig) requestFriendHandler(w http.ResponseWriter, r *http.Request) {
	type addFriendRes struct {
		RelationID string    `json:"relation_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		UserID     string    `json:"user_id"`
		ReceiverID string    `json:"receiver_id"`
		Status     string    `json:"status"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid ID", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, apiCfg.SecretJWT)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid ID", err)
		return
	}

	user, err := apiCfg.Db.GetUserByID(r.Context(), userID)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid ID", err)
		return
	}

	receiverLogin := r.PathValue("receiverLogin")
	receiver, err := apiCfg.Db.GetUserByLogin(r.Context(), receiverLogin)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "User not found", err)
		return
	}

	newFriendshipParams := database.CreateFriendshipParams{
		UserID:     user.ID,
		ReceiverID: receiver.ID,
	}
	relation, err := apiCfg.Db.CreateFriendship(r.Context(), newFriendshipParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	successResponse(w, http.StatusCreated, addFriendRes{
		CreatedAt:  relation.CreatedAt,
		UpdatedAt:  relation.UpdatedAt,
		UserID:     relation.UserID.String(),
		ReceiverID: relation.ReceiverID.String(),
		Status:     string(relation.Relationship),
	})
}
