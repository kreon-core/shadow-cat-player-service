package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/kreon-core/shadow-cat-common/appc"
	"github.com/kreon-core/shadow-cat-common/resc"

	"sc-player-service/helper"
	"sc-player-service/infrastructure/external"
	"sc-player-service/model/api/dto"
	"sc-player-service/model/api/response"
)

type Auth struct {
	AuthClient *external.HTTPClient
}

func NewAuthMiddleware(authClient *external.HTTPClient) *Auth {
	return &Auth{
		AuthClient: authClient,
	}
}

func (m *Auth) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		headers := map[string]string{}
		if token := r.Header.Get("Authorization"); token != "" {
			headers["Authorization"] = token
		}
		code, resp, data, err := external.CallAPI[dto.VerifiedTokenData](
			r.Context(), m.AuthClient,
			http.MethodGet, m.AuthClient.Paths["verify-user"],
			nil, headers)
		if err != nil {
			resc.JSON(w, http.StatusInternalServerError, &response.Resp{
				ReturnCode:    appc.EExternalServiceError,
				ReturnMessage: appc.Message(appc.EExternalServiceError),
			})
			return
		}

		if resp.ReturnCode != appc.Success {
			resc.JSON(ww, code, resp)
			return
		}

		ctx := context.WithValue(r.Context(), helper.PlayerIDContextKey, data.UserID)

		next.ServeHTTP(ww, r.WithContext(ctx))
	})
}
