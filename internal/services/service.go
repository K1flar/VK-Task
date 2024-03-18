package services

import (
	"film_library/internal/config"
	"film_library/internal/domains"
	"film_library/internal/repositories/postgres"
	"film_library/internal/services/actorservice"
	"film_library/internal/services/filmservice"
	userservicce "film_library/internal/services/userservice"
	"film_library/pkg/pagination"
	"log/slog"
	"time"
)

type UserService interface {
	CreateUser(user domains.User) (string, error)
	Login(login, password string) (string, error)
}

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

type ActorService interface {
	CreateActor(actor domains.Actor) error
	AddActorsToFilm(filmID uint32, actorsID []uint32) error
	UpdateActorFullName(id uint32, fullName string) error
	UpdateActorGender(id uint32, gender domains.Gender) error
	UpdateActorBirthday(id uint32, birthday time.Time) error
	UpdateActor(id uint32, actor domains.Actor) error
	DeleteActor(id uint32) error
	DeleteActorFromFilm(actorID uint32, filmID uint32) error
	GetActorsWithFilms(filter *pagination.ActorsFilter) ([]*domains.ActorWithFilms, error)
}

type Service struct {
	UserService
	FilmService
	ActorService
}

type IService interface {
	UserService
	FilmService
	ActorService
}

func New(repo postgres.IRepository, log *slog.Logger, cfg *config.Config) IService {
	userService := userservicce.New(repo, log, cfg)
	actorService := actorservice.New(repo, log)
	filmservice := filmservice.New(repo, actorService, log, cfg)
	return &Service{
		userService,
		filmservice,
		actorService,
	}
}
