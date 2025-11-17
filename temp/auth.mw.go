package temp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	tul "github.com/kreon-core/shadow-cat-common"

	"sc-player-service/infrastructure/config"
	"sc-player-service/model/api/response"
)

type Auth struct {
	Secrets *config.Secrets
}

func NewAuthMiddleware(secrets *config.Secrets) *Auth {
	return &Auth{
		Secrets: secrets,
	}
}

// Context keys.
type contextKey string

const (
	ContextIsAuthenticated contextKey = "is_authenticated"
	ContextUserID          contextKey = "user_id"
	ContextRole            contextKey = "role"
)

func (m *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val := r.Context().Value(ContextIsAuthenticated); val != nil {
			if bl, ok := val.(bool); ok && bl {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, tul.EMissingAuthorizationHeader)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, claims, err := ParseJWTAccessToken(tokenString, []byte(m.Secrets.JWTSecretKey))
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				writeJSON(w, http.StatusUnauthorized, tul.EExpiredAccessToken)
			} else {
				writeJSON(w, http.StatusUnauthorized, tul.EInvalidAccessToken)
			}
			return
		}
		if !token.Valid {
			writeJSON(w, http.StatusUnauthorized, tul.EInvalidAccessToken)
			return
		}

		ctx := context.WithValue(r.Context(), ContextIsAuthenticated, true)
		ctx = context.WithValue(ctx, ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextRole, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSON(w http.ResponseWriter, status, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := &response.Resp{
		ReturnCode:    code,
		ReturnMessage: tul.Message(code),
	}
	json.NewEncoder(w).Encode(resp)
}
