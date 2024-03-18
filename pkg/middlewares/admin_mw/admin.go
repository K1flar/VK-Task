package adminmw

import (
	"film_library/internal/domains"
	"film_library/internal/handlers/response"
	"film_library/pkg/middlewares/auth"
	"log/slog"
	"net/http"
)

const AdminRole = "admin"

func New(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(auth.UserKey("user")).(domains.User)
			if !ok || user.Role != AdminRole {
				response.JSONError(w, http.StatusForbidden, "forbidden", log)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
