package handlers

import (
	"film_library/internal/handlers/actorhandler"
	"film_library/internal/handlers/filmhandler"
	"film_library/internal/handlers/userhandler"
	"film_library/internal/services"
	"log/slog"
)

type Handler struct {
	*userhandler.UserHandler
	*actorhandler.ActorHandler
	*filmhandler.FilmHandler
}

func New(service services.IService, log *slog.Logger) *Handler {
	return &Handler{
		userhandler.New(service, log),
		actorhandler.New(service, log),
		filmhandler.New(service, log),
	}
}
