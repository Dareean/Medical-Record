package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/JinXVIII/BE-Medical-Record/internal/domain"
	"github.com/JinXVIII/BE-Medical-Record/pkg/helper"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.SendJSON(w, http.StatusUnauthorized, domain.Response{Message: "Token tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			helper.SendJSON(w, http.StatusUnauthorized, domain.Response{Message: "Format token salah"})
			return
		}

		// Get JWT secret from environment
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			helper.SendJSON(w, http.StatusInternalServerError, domain.Response{Message: "Server configuration error"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			helper.SendJSON(w, http.StatusUnauthorized, domain.Response{Message: "Token tidak valid"})
			return
		}

		// Add user info to context
		if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
			userInfo := make(map[string]interface{})
			userInfo["user_id"] = mapClaims["user_id"]
			userInfo["email"] = mapClaims["email"]
			userInfo["role"] = mapClaims["role"]

			// Add to request context
			ctx := context.WithValue(r.Context(), "user", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		helper.SendJSON(w, http.StatusUnauthorized, domain.Response{Message: "Token tidak valid"})
	})
}
