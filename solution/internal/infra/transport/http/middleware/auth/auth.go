package mwauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
)

var secretKey = []byte(os.Getenv("jwt_secret"))

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerString := r.Header.Get("Authorization")
		if headerString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(utils.MakeErr("missing authorization header"))
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		parts := strings.Split(headerString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(utils.MakeErr("invalid Authorization header format"))
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}
		tokenString := parts[1]

		token, err := parseToken(tokenString)
		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(utils.MakeErr("invalid token"))
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		var ctx context.Context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx = context.WithValue(r.Context(), "username", claims["username"])
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}
