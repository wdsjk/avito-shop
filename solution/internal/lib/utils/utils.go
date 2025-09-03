package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
)

func MakeErr(msg string) handlers_dto.ErrorResponse {
	return handlers_dto.ErrorResponse{Errors: msg}
}

func GenerateJWT(username string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(secret)
}
