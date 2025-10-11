package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/MedrekIT/message-api/internal/auth"
	"github.com/MedrekIT/message-api/internal/database"
	"github.com/google/uuid"
)

type groupRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	CreatorID uuid.UUID `json:"creator_id"`
	GroupType string    `json:"group_type"`
}

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
		errorResponse(w, http.StatusBadRequest, "Invalid group ID", err)
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

func (apiCfg *ApiConfig) getInvitationGroupHandler(w http.ResponseWriter, r *http.Request) {
	type groupFromInvitationRes struct {
		GroupID string `json:"group_id"`
	}

	invitationToken := r.PathValue("invitationToken")
	invitationLink, err := apiCfg.Db.GetInvitation(r.Context(), invitationToken)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	successResponse(w, http.StatusOK, groupFromInvitationRes{
		GroupID: invitationLink.OfGroupID.String(),
	})
}

func (apiCfg *ApiConfig) getGroupHandler(w http.ResponseWriter, r *http.Request) {
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
		errorResponse(w, http.StatusBadRequest, "Invalid group ID", err)
		return
	}
	group, err := apiCfg.Db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Invalid group ID", err)
		return
	}

	getMemberParams := database.GetMemberParams{
		UserID:    user.ID,
		OfGroupID: group.ID,
	}
	_, err = apiCfg.Db.GetMember(r.Context(), getMemberParams)
	if err != nil && group.GroupType != "public" {
		errorResponse(w, http.StatusForbidden, "You cannot see that data", err)
		return
	}

	successResponse(w, http.StatusOK, groupRes{
		ID:        group.ID,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		Name:      group.Name,
		CreatorID: group.CreatorID,
		GroupType: string(group.GroupType),
	})
}

func (apiCfg *ApiConfig) getMembersHandler(w http.ResponseWriter, r *http.Request) {
	type membersInfo struct {
		UserID     string `json:"user_id"`
		UserLogin  string `json:"user_login"`
		MemberType string `json:"member_type"`
	}
	type membersRes struct {
		Members []membersInfo `json:"members"`
	}

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
		errorResponse(w, http.StatusBadRequest, "Invalid group ID", err)
		return
	}
	group, err := apiCfg.Db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Invalid group ID", err)
		return
	}

	getMemberParams := database.GetMemberParams{
		UserID:    user.ID,
		OfGroupID: group.ID,
	}
	_, err = apiCfg.Db.GetMember(r.Context(), getMemberParams)
	if err != nil {
		errorResponse(w, http.StatusForbidden, "You cannot see that data", err)
		return
	}

	members, err := apiCfg.Db.GetMembers(r.Context(), group.ID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	var membersOfGroup []membersInfo
	for _, member := range members {
		membersOfGroup = append(membersOfGroup, membersInfo{
			UserID:     member.UserID.String(),
			UserLogin:  member.Login,
			MemberType: string(member.MemberType),
		})
	}

	successResponse(w, http.StatusOK, membersRes{
		Members: membersOfGroup,
	})
}

func (apiCfg *ApiConfig) joinGroupHandler(w http.ResponseWriter, r *http.Request) {
	type joinGroupReq struct {
		InvitationLink string `json:"invitation_link"` // optional if group is public
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	var reqData joinGroupReq
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

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
		errorResponse(w, http.StatusBadRequest, "Invalid group ID", err)
		return
	}
	group, err := apiCfg.Db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Invalid group ID", err)
		return
	}

	if group.GroupType != "public" {
		var invitation database.InvitationLink
		if reqData.InvitationLink == "" || group.GroupType == "private" {
			errorResponse(w, http.StatusUnauthorized, "Unauthorized operation", nil)
			return
		} else if reqData.InvitationLink != "" {
			invitation, err = apiCfg.Db.GetInvitation(r.Context(), reqData.InvitationLink)
			if err != nil {
				errorResponse(w, http.StatusBadRequest, "Invalid inviation link", err)
				return
			}
		}

		if reqData.InvitationLink != "" && invitation.OfGroupID != group.ID {
			errorResponse(w, http.StatusBadRequest, "Invalid inviation link", err)
			return
		}
	}

	newMemberParams := database.AddMemberParams{
		OfGroupID: group.ID,
		UserID:    user.ID,
	}
	_, err = apiCfg.Db.AddMember(r.Context(), newMemberParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	successResponse(w, http.StatusCreated, groupRes{
		ID:        group.ID,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		Name:      group.Name,
		CreatorID: group.CreatorID,
		GroupType: string(group.GroupType),
	})
}

func (apiCfg *ApiConfig) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	type createGroupReq struct {
		Name      string `json:"name"`
		GroupType string `json:"group_type"` // optional ('public', 'invite_only', 'private') set to 'invite_only' by default
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	var reqData createGroupReq
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if reqData.GroupType == "" {
		reqData.GroupType = "invite_only"
	}

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

	newGroupParams := database.CreateGroupParams{
		ID:        uuid.New(),
		Name:      reqData.Name,
		CreatorID: user.ID,
		GroupType: database.GroupType(reqData.GroupType),
	}
	group, err := apiCfg.Db.CreateGroup(r.Context(), newGroupParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	newMemberParams := database.AddMemberParams{
		OfGroupID: group.ID,
		UserID:    user.ID,
	}
	member, err := apiCfg.Db.AddMember(r.Context(), newMemberParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		err := apiCfg.Db.DeleteGroup(r.Context(), group.ID)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
		return
	}

	newPermissionsParams := database.ChangePermissionsParams{
		OfGroupID:  member.OfGroupID,
		UserID:     member.UserID,
		MemberType: "admin",
	}
	err = apiCfg.Db.ChangePermissions(r.Context(), newPermissionsParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		err := apiCfg.Db.DeleteGroup(r.Context(), group.ID)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
		return
	}

	successResponse(w, http.StatusCreated, groupRes{
		ID:        group.ID,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		Name:      group.Name,
		CreatorID: group.CreatorID,
		GroupType: string(group.GroupType),
	})
}
