package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MedrekIT/message-api/internal/auth"
	"github.com/MedrekIT/message-api/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type successRes struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Login     string    `json:"login"`
		Email     string    `json:"email"`
	}
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", fmt.Errorf("error while reading request body - %w\n", err))
		return
	}

	var reqData loginReq
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", fmt.Errorf("error while decoding request body - %w\n", err))
		return
	}

	user, err := apiCfg.Db.GetUserByEmail(r.Context(), reqData.Email)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "User with given e-mail does not exist", err)
		return
	}

	if isValid, err := auth.CheckPasswordHash(reqData.Password, user.Password); !isValid || err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	successResponse(w, 201, successRes{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Login:     user.Login,
		Email:     user.Email,
	})
}

func (apiCfg *ApiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {
	type successRes struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Login     string    `json:"login"`
		Email     string    `json:"email"`
	}
	type registerReq struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", fmt.Errorf("error while reading request body - %w\n", err))
		return
	}

	var reqData registerReq
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request", fmt.Errorf("error while decoding request body - %w\n", err))
		return
	}

	hashedPassword, err := auth.HashPassword(reqData.Password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	newUserParams := database.CreateUserParams{
		ID:       uuid.New(),
		Login:    reqData.Login,
		Email:    reqData.Email,
		Password: hashedPassword,
	}
	user, err := apiCfg.Db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		errorResponse(w, http.StatusConflict, "User with given login or e-mail already exists", err)
		return
	}

	successResponse(w, 201, successRes{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Login:     user.Login,
		Email:     user.Email,
	})
}
