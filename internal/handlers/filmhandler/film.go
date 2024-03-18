package filmhandler

import (
	"encoding/json"
	"errors"
	"film_library/internal/domains"
	"film_library/internal/handlers/response"
	"film_library/internal/repositories/postgres/filmrepo"
	"film_library/pkg/pagination"
	"film_library/pkg/validation"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type FilmService interface {
	CreateFilm(film domains.Film, actors []uint32) (uint32, error)
	UpdateFilmName(id uint32, name string) error
	UpdateFilmDescription(id uint32, descrtion string) error
	UpdateFilmReleaseDate(id uint32, releaseDate time.Time) error
	UpdateFilmRating(id uint32, rating int) error
	UpdateFilm(id uint32, film domains.Film) error
	DeleteFilm(id uint32) error
	GetFilms(filter *pagination.FilmFilter) ([]*domains.Film, error)
}

type FilmHandler struct {
	service FilmService
	log     *slog.Logger
}

func New(service FilmService, log *slog.Logger) *FilmHandler {
	return &FilmHandler{
		service: service,
		log:     log,
	}
}

type InputCreateFilm struct {
	Film     domains.Film `json:"film"`
	ActorsID []uint32     `json:"actorsID"`
}

// @Summary Create film
// @Tags film
// @Description create film
// @ID create-film
// @Accept  json
// @Produce  json
// @Param input body InputCreateFilm true "film and actors info"
// @Success 200 {object} integer
// @Success 400 {object} response.ErrorsReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film [post]
func (h *FilmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}
	defer r.Body.Close()

	input := InputCreateFilm{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	film := input.Film
	actors := input.ActorsID

	id, err := h.service.CreateFilm(film, actors)
	if err != nil {
		if err, ok := err.(*validation.ValidateError); ok {
			response.JSONErrors(w, http.StatusBadRequest, err.ToArrayErrors(), h.log)
			return
		}
		if errors.Is(err, filmrepo.ErrAlreadyExists) {
			response.JSONError(w, http.StatusBadRequest, "film already exists", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"id": id,
	}, h.log)
}

// @Summary Get films
// @Tags film
// @Description get films
// @ID get-films
// @Accept  json
// @Produce  json
// @Param page query integer false "page number"
// @Param size query integer false "page size"
// @Param film query string false "film name contains"
// @Param actor query string false "actor full name contains"
// @Param sort query string false "films order by"
// @Success 200 {object} []domains.Film
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/films [get]
func (h *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	filter := pagination.NewFilmFilterFromRequest(r)

	actorsWithFilms, err := h.service.GetFilms(filter)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSON(w, http.StatusOK, actorsWithFilms, h.log)
}

// @Summary Update film name
// @Tags film
// @Description update film name
// @ID update-name
// @Accept  json
// @Produce  json
// @Param id path integer true "film id"
// @Param name path string true "film name"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film/name/{id}/{name} [put]
func (h *FilmHandler) UpdateFilmName(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	name := r.PathValue("name")

	err = h.service.UpdateFilmName(uint32(id), name)
	if err != nil {
		if errors.Is(err, filmrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "film not found", h.log)
			return
		}
		if errors.Is(err, filmrepo.ErrInvalidNameLength) {
			response.JSONError(w, http.StatusBadRequest, "invalid film name length", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type InputDescription struct {
	Description string `json:"description"`
}

// @Summary Update film description
// @Tags film
// @Description update film description
// @ID update-description
// @Accept  json
// @Produce  json
// @Param id path integer true "film id"
// @Param description body InputDescription true "actor gender"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film/description/{id} [put]
func (h *FilmHandler) UpdateFilmDescription(w http.ResponseWriter, r *http.Request) {
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

	description := InputDescription{}
	err = json.Unmarshal(b, &description)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.UpdateFilmDescription(uint32(id), description.Description)
	if err != nil {
		if errors.Is(err, filmrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "film not found", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update film release date
// @Tags film
// @Description update film release date
// @ID update-releaseDate
// @Accept  json
// @Produce  json
// @Param id path integer true "film id"
// @Param date path string true "film release date" format(2006-01-02)
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film/date/{id}/{date} [put]
func (h *FilmHandler) UpdateFilmReleaseDate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	releaseDate, err := time.Parse(time.DateOnly, r.PathValue("date"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid date", h.log)
		return
	}

	err = h.service.UpdateFilmReleaseDate(uint32(id), releaseDate)
	if err != nil {
		if errors.Is(err, filmrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update film rating
// @Tags film
// @Description update film rating
// @ID update-rating
// @Accept  json
// @Produce  json
// @Param id path integer true "film id"
// @Param rating path integer true "film rating"
// @Success 200
// @Failure 400 {object} response.ErrorsReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film/{id}/{rating} [put]
func (h *FilmHandler) UpdateFilmRating(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	rating, err := time.Parse(time.DateOnly, r.PathValue("rating"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid date", h.log)
		return
	}

	err = h.service.UpdateFilmReleaseDate(uint32(id), rating)
	if err != nil {
		if errors.Is(err, filmrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "actor not found", h.log)
			return
		}
		if errors.Is(err, filmrepo.ErrInvalidRating) {
			response.JSONError(w, http.StatusInternalServerError, "invalid rating", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update film
// @Tags film
// @Description update film
// @ID update-film
// @Accept  json
// @Produce  json
// @Param id path integer true "film id"
// @Param input body domains.Film true "film info"
// @Success 200
// @Failure 400 {object} response.ErrorsReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film/{id} [put]
func (h *FilmHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
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

	film := domains.Film{}
	err = json.Unmarshal(b, &film)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request kurwa", h.log)
		return
	}

	err = h.service.UpdateFilm(uint32(id), film)
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

// @Summary Delete film
// @Tags film
// @Description delete film by id
// @ID delete-film
// @Accept  json
// @Produce  json
// @Param id path integer true "film id"
// @Success 200
// @Failure 400 {object} response.ErrorReponse
// @Failure 500 {object} response.ErrorReponse
// @Security ApiKeyAuth
// @Router /api/film/{id} [delete]
func (h *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	err = h.service.DeleteFilm(uint32(id))
	if err != nil {
		if errors.Is(err, filmrepo.ErrNotFound) {
			response.JSONError(w, http.StatusBadRequest, "film not found", h.log)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}
