package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)

	return hex.EncodeToString(key)
}

func CreateJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "messageAPI",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error while signing token with secret key - %w\n", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing JWT token - %w\n", err)
	}

	stringID, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.UUID{}, err
	}

	if issuer != "messageAPI" {
		return uuid.UUID{}, fmt.Errorf("invalid issuer\n")
	}

	userID, err := uuid.Parse(stringID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user ID - %w\n", err)
	}

	return userID, nil
}

func GetBearerToken(header http.Header) (string, error) {
	tokenString := header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("authorization header is empty\n")
	}

	tokenArr := strings.Split(tokenString, " ")
	if len(tokenArr) != 2 || tokenArr[0] != "Bearer" {
		return "", fmt.Errorf("invalid token body\n")
	}

	return tokenArr[1], nil
}
