package middleware

import (
	"context"
	"net/http"
	"strings"
	"task-manager-backend/data/common"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationKey = "Authorization"
)

type JWTAuthorizationMiddleware struct {
	secret 		string
}

func NewJWTAuthorizationMiddleware(secret string) *JWTAuthorizationMiddleware {
	return &JWTAuthorizationMiddleware{
		secret: secret,
	}
}

func (m *JWTAuthorizationMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationKey)

		// Check for the presence of valid JWT token in the header
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &common.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), common.USER_ID_KEY, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}
