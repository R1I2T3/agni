package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func ParseJWT(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnvConfig().AdminEnvConfig.JWT_Secret), nil
	})
	return token, claims, err
}
