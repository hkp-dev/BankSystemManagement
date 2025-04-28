package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "userId"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")

			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, `{"error": "Missing auth token"}`, http.StatusUnauthorized)
				return
			}

			// Parse token với JWT v5
			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				// Kiểm tra signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			// Xử lý lỗi token
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					http.Error(w, `{"error": "Auth token expired"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error": "Invalid auth token"}`, http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, `{"error": "Invalid auth token"}`, http.StatusUnauthorized)
				return
			}

			// Lấy userId từ claims
			userID, ok := (*claims)["userId"].(string)
			if !ok || userID == "" {
				http.Error(w, `{"error": "Invalid user ID in token"}`, http.StatusUnauthorized)
				return
			}

			// Thêm userId vào context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
