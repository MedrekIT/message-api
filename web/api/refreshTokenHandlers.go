package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MedrekIT/message-api/internal/auth"
)

func (apiCfg *ApiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type successRes struct {
		Token string `json:"token"`
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	refreshToken, err := apiCfg.Db.GetRefreshToken(r.Context(), reqToken)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("error while getting refresh token from the database - %w\n", err))
		return
	}

	newJWT, err := auth.CreateJWT(refreshToken.UserID, apiCfg.SecretJWT, time.Hour)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	successResponse(w, http.StatusOK, successRes{
		Token: newJWT,
	})
}

func (apiCfg *ApiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	err = apiCfg.Db.RevokeRefreshToken(r.Context(), reqToken)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("error while getting refresh token from the database - %w\n", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
