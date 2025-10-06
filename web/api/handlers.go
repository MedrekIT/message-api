package api

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"

	"github.com/MedrekIT/message-api/internal/auth"
	"github.com/MedrekIT/message-api/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {
	type successRes struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	type registerReq struct {
		Login string `json:"login"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var reqData registerReq
	if err := json.Unmarshal(r.Body, &reqData); err != nil {
		errorResponse(w, http.BadRequest, "Invalid request", fmt.Sprintf("error while decoding request body - %w\n", err))
		return
	}

	hashedPassword, err := auth.HashPassword(reqData.Password)
	if err != nil {
		errorResponse(w, http.InternalServerError, "Something went wrong", err)
		return
	}

	newUserParams := database.CreateUserParams{
		ID: uuid.New(),
		Login: reqData.Login,
		Email: reqData.Email,
		Password: hashedPassword,
	}
	user, err := apiCfg.Db.CreateUser(req.Context(), newUserParams)
	if err != nil {
		errorResponse(w, http.Conflict, "User with given login or e-mail already exists", err)
		return
	}

	successResponse(resWriter, 201, successRes{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Login: user.Login,
		Email: user.Email,
	})
}
