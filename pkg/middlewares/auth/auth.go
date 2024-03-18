package auth

import (
	"context"
	"encoding/json"
	"film_library/internal/config"
	"film_library/internal/domains"
	"film_library/internal/handlers/response"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserKey string

func New(log *slog.Logger, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authParts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(authParts) < 2 {
				response.JSONError(w, http.StatusUnauthorized, "unauthorized", log)
				return
			}

			inToken := authParts[1]
			token, err := jwt.Parse(inToken, func(t *jwt.Token) (interface{}, error) {
				if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != "HS256" {
					return nil, fmt.Errorf("bad sign method")
				}
				return []byte(cfg.Server.Secret), nil
			})
			if err != nil || !token.Valid {
				response.JSONError(w, http.StatusUnauthorized, "bad token", log)
				return
			}

			payload, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.JSONError(w, http.StatusUnauthorized, "no payload", log)
				return
			}

			userStr, err := json.Marshal(payload)
			if err != nil {
				response.JSONError(w, http.StatusInternalServerError, "unknown error", log)
				return
			}

			var userStruct domains.User
			err = json.Unmarshal(userStr, &userStruct)
			if err != nil {
				response.JSONError(w, http.StatusInternalServerError, "unknown error", log)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey("user"), userStruct)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
