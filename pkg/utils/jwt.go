package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/r1i2t3/agni/pkg/config"
)

func GenerateAdminJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"admin":    true,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	envConfig := config.GetEnvConfig()
	return token.SignedString([]byte(envConfig.AdminEnvConfig.JWT_Secret))
}

// GenerateApplicationJWT generates a JWT for application clients
func GenerateApplicationJWT(applicationID uuid.UUID, userID string) (string, error) {
	claims := jwt.MapClaims{
		"application_id": applicationID.String(),
		"user_id":        userID,
		"type":           "application",
		"exp":            time.Now().Add(24 * time.Hour).Unix(),
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	envConfig := config.GetEnvConfig()
	return token.SignedString([]byte(envConfig.AdminEnvConfig.JWT_Secret))
}

func ParseJWT(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnvConfig().AdminEnvConfig.JWT_Secret), nil
	})
	return token, claims, err
}

// ValidateApplicationJWT validates and extracts claims from application JWT
func ValidateApplicationJWT(tokenStr string) (applicationID uuid.UUID, userID string, err error) {
	token, claims, err := ParseJWT(tokenStr)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, "", fmt.Errorf("token is not valid")
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "application" {
		return uuid.Nil, "", fmt.Errorf("invalid token type")
	}

	// Extract application ID
	appIDStr, ok := claims["application_id"].(string)
	if !ok {
		return uuid.Nil, "", fmt.Errorf("missing application_id in token")
	}

	applicationID, err = uuid.Parse(appIDStr)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid application_id format")
	}

	// Extract user ID
	userID, ok = claims["user_id"].(string)
	if !ok {
		return uuid.Nil, "", fmt.Errorf("missing user_id in token")
	}

	return applicationID, userID, nil
}
