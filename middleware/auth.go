package middleware

import (
	"context"
	"net/http"
	"project-bioskop/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const KeyUserID contextKey = "userID"
const KeyUserRole contextKey = "userRole"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			ctx := context.WithValue(r.Context(), KeyUserID, claims["user_id"])
			ctx = context.WithValue(ctx, KeyUserRole, claims["role"])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}