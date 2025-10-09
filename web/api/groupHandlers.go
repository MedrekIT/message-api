package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/MedrekIT/message-api/internal/auth"
	"github.com/MedrekIT/message-api/internal/database"
	"github.com/google/uuid"
)

func CreateInvitationKey() string {
	key := make([]byte, 8)
	rand.Read(key)

	return hex.EncodeToString(key)
}

func (apiCfg *ApiConfig) createInvitationHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, apiCfg.SecretJWT)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	user, err := apiCfg.Db.GetUserByID(r.Context(), userID)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	groupID, err := uuid.Parse(r.PathValue("groupID"))
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Invalid group ID", err)
		return
	}
	group, err := apiCfg.Db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Invalid group ID", err)
		return
	}

	if group.CreatorID != user.ID || group.GroupType == "private" {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized operation", nil)
		return
	}
	invitationKey := CreateInvitationKey()

	newInvitationLinkParams := database.CreateInvitationLinkParams{
		Token:     invitationKey,
		OfGroupID: group.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	invitationLink, err := apiCfg.Db.CreateInvitationLink(r.Context(), newInvitationLinkParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	successResponse(w, http.StatusCreated, invitationLink)
}
