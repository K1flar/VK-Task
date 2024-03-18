package postgres

import (
	"database/sql"
	"film_library/internal/config"
	"film_library/internal/domains"
	"film_library/internal/repositories/postgres/actorrepo"
	"film_library/internal/repositories/postgres/filmrepo"
	"film_library/internal/repositories/postgres/userrepo"
	"film_library/pkg/pagination"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type UserRepo interface {
	AddUser(user domains.User) error
	GetUserByLoign(login string) (*domains.User, error)
}

type ActorRepo interface {
	AddActor(actor domains.Actor) error
	AddActorsToFilm(filmID uint32, actorsID []uint32) error
	UpdateActorFullName(id uint32, fullName string) error
	UpdateActorGender(id uint32, gender string) error
	UpdateActorBirthday(id uint32, birthday time.Time) error
	UpdateActor(id uint32, actor domains.Actor) error
	DeleteActor(id uint32) error
	DeleteActorFromFilm(actorID uint32, filmID uint32) error
	GetActorsWithFilms(filter *pagination.ActorsFilter) ([]*domains.ActorWithFilms, error)
}

type FilmRepo interface {
	AddFilm(film domains.Film) (uint32, error)
	UpdateFilmName(id uint32, name string) error
	UpdateFilmDescription(id uint32, descrtion string) error
	UpdateFilmReleaseDate(id uint32, releaseDate time.Time) error
	UpdateFilmRating(id uint32, rating int) error
	UpdateFilm(id uint32, film domains.Film) error
	DeleteFilm(id uint32) error
	GetFilms(filter *pagination.FilmFilter) ([]*domains.Film, error)
}

type IRepository interface {
	UserRepo
	ActorRepo
	FilmRepo
}

type Repository struct {
	UserRepo
	ActorRepo
	FilmRepo
}

func New(cfg *config.DataBase) (IRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{
		userrepo.NewUserRepository(db),
		actorrepo.NewActorRepository(db),
		filmrepo.NewFilmRepository(db),
	}, nil
}
