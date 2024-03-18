package actorhandler

import (
	"encoding/json"
	"errors"
	"film_library/internal/domains"
	"film_library/internal/handlers/response"
	"film_library/internal/repositories/postgres/actorrepo"
	"film_library/internal/services/actorservice"
	"film_library/pkg/pagination"
	"film_library/pkg/validation"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ActorService interface {
	CreateActor(actor domains.Actor) error
	AddActorsToFilm(filmID uint32, actors []uint32) error
	UpdateActorFullName(id uint32, fullName string) error
	UpdateActorGender(id uint32, gender domains.Gender) error
	UpdateActorBirthday(id uint32, birthday time.Time) error
	UpdateActor(id uint32, actor domains.Actor) error
	DeleteActor(id uint32) error
	DeleteActorFromFilm(actorID uint32, filmID uint32) error
	GetActorsWithFilms(filter *pagination.ActorsFilter) ([]*domains.ActorWithFilms, error)
}

type ActorHandler struct {
	service ActorService
	log     *slog.Logger
}

func New(service ActorService, log *slog.Logger) *ActorHandler {
	return &ActorHandler{
		service: service,
		log:     log,
	}
}

// @Summary Get actors with films
// @Tags actor
// @Description get actors with films
// @ID get-actors
// @Accept  json
// @Produce  json
// @Param page query integer false "page number"
// @Param size query integer false "page size"
// @Param actor query string false "full name contains"
// @Success 200 {object} []domains.ActorWithFilms
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actors [get]
func (h *ActorHandler) GetActorsWithFilms(w http.ResponseWriter, r *http.Request) {
	filter := pagination.NewActorFilterFromRequest(r)

	actorsWithFilms, err := h.service.GetActorsWithFilms(filter)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSON(w, http.StatusOK, actorsWithFilms, h.log)
}

// @Summary Create actor
// @Tags actor
// @Description create actor
// @ID create-actor
// @Accept  json
// @Produce  json
// @Param input body domains.Actor true "actor info"
// @Success 200
// @Failure 400 {object} response.ErrorsReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor [post]
func (h *ActorHandler) CreateActor(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}
	defer r.Body.Close()

	var actor domains.Actor
	err = json.Unmarshal(b, &actor)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.CreateActor(actor)
	if err != nil {
		if err, ok := err.(*validation.ValidateError); ok {
			response.JSONErrors(w, http.StatusBadRequest, err.ToArrayErrors(), h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Add actors to film
// @Tags actor
// @Description add actors to film
// @ID add-actors
// @Accept  json
// @Produce  json
// @Param filmID path integer true "film id"
// @Param input body []uint32 true "actors id"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actors/{filmID} [post]
func (h *ActorHandler) AddActorsToFilm(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}
	defer r.Body.Close()

	filmID, err := strconv.Atoi(r.PathValue("filmID"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	var actorsID []uint32
	err = json.Unmarshal(b, &actorsID)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.AddActorsToFilm(uint32(filmID), actorsID)
	if err != nil {
		if errors.Is(err, actorrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		if errors.Is(err, actorrepo.ErrUniqueActors) {
			response.JSONError(w, http.StatusBadRequest, "actors must be unique", h.log)
			return
		}

		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update actor full name
// @Tags actor
// @Description update actor full name
// @ID update-fullname
// @Accept  json
// @Produce  json
// @Param id path integer true "actor id"
// @Param name path string true "actor full name"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor/name/{id}/{name} [put]
func (h *ActorHandler) UpdateActorFullName(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	fullName := r.PathValue("name")

	err = h.service.UpdateActorFullName(uint32(id), fullName)
	if err != nil {
		if errors.Is(err, actorrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		if errors.Is(err, actorservice.ErrInvalidFullName) {
			response.JSONError(w, http.StatusBadRequest, "full name must be at least 1 letter long", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update actor gender
// @Tags actor
// @Description update actor gender
// @ID update-gender
// @Accept  json
// @Produce  json
// @Param id path integer true "actor id"
// @Param gender path string true "actor gender"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor/gender/{id}/{gender} [put]
func (h *ActorHandler) UpdateActorGender(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	gender := r.PathValue("gender")

	err = h.service.UpdateActorGender(uint32(id), domains.Gender(gender))
	if err != nil {
		if errors.Is(err, actorrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		if errors.Is(err, actorservice.ErrInvalidGender) {
			response.JSONError(w, http.StatusBadRequest, "gender must be male or female", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update actor birthday
// @Tags actor
// @Description update actor birthday
// @ID update-birthday
// @Accept  json
// @Produce  json
// @Param id path integer true "actor id"
// @Param birthday path string true "actor birthday" format(2006-01-02)
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor/birthday/{id}/{birthday} [put]
func (h *ActorHandler) UpdateActorBirthday(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	birthday, err := time.Parse(time.DateOnly, r.PathValue("birthday"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid date", h.log)
		return
	}

	err = h.service.UpdateActorBirthday(uint32(id), birthday)
	if err != nil {
		if errors.Is(err, actorrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update actor
// @Tags actor
// @Description update actor
// @ID update-actor
// @Accept  json
// @Produce  json
// @Param id path integer true "actor id"
// @Param input body domains.Actor true "actor info"
// @Success 200
// @Failure 400 {object} response.ErrorsReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor/{id} [put]
func (h *ActorHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}
	defer r.Body.Close()

	actor := domains.Actor{}
	err = json.Unmarshal(b, &actor)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.UpdateActor(uint32(id), actor)
	if err != nil {
		if err, ok := err.(*validation.ValidateError); ok {
			response.JSONErrors(w, http.StatusBadRequest, err.ToArrayErrors(), h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete actor
// @Tags actor
// @Description delete actor by id
// @ID delete-actor
// @Accept  json
// @Produce  json
// @Param id path integer true "actor id"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor/{id} [delete]
func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.DeleteActor(uint32(id))
	if err != nil {
		if errors.Is(err, actorrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete actor from film
// @Tags actor
// @Description delete actor from film
// @ID delete-actor-from-film
// @Accept  json
// @Produce  json
// @Param id path integer true "actor id"
// @Param filmID path integer true "film id"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/actor/{id}/{filmID} [delete]
func (h *ActorHandler) DeleteActorFromFilm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}
	filmID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.DeleteActorFromFilm(uint32(id), uint32(filmID))
	if err != nil {
		if errors.Is(err, actorrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}
