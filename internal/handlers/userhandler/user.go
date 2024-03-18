package userhandler

import (
	"encoding/json"
	"errors"
	"film_library/internal/domains"
	"film_library/internal/handlers/response"
	"film_library/internal/repositories/postgres/userrepo"
	"film_library/internal/services/userservice"
	"film_library/pkg/validation"
	"io"
	"log/slog"
	"net/http"
)

type UserService interface {
	CreateUser(user domains.User) (string, error)
	Login(login, password string) (string, error)
}

type UserHandler struct {
	service UserService
	log     *slog.Logger
}

func New(service UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log,
	}
}

// @Summary Create user
// @Tags user
// @Description create user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body domains.User true "user info"
// @Success 200 {object} string ""
// @Failure 400 {object} response.ErrorReponse
// @Failure 400 {object} response.ErrorsReponse
// @Failure 409 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Router /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}
	defer r.Body.Close()

	var user domains.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	token, err := h.service.CreateUser(user)
	if err != nil {
		if err, ok := err.(*validation.ValidateError); ok {
			response.JSONErrors(w, http.StatusBadRequest, err.ToArrayErrors(), h.log)
			return
		}
		if errors.Is(err, userrepo.ErrAlreadyExists) {
			response.JSONError(w, http.StatusConflict, "user already exists", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"token": token,
	}, h.log)
}

// @Summary Login user
// @Tags user
// @Description login user
// @ID login
// @Accept  json
// @Produce  json
// @Param input body domains.User true "user info"
// @Success 200 {object} string ""
// @Failure 400 {object} response.ErrorReponse
// @Failure 401 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Router /api/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error kurwa", h.log)
		return
	}
	defer r.Body.Close()

	var user domains.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	token, err := h.service.Login(user.Login, user.Password)
	if err != nil {
		if errors.Is(err, userservice.ErrNotFound) || errors.Is(err, userservice.ErrInvalidPassword) {
			response.JSONError(w, http.StatusUnauthorized, "invalid login or password", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"token": token,
	}, h.log)
}
