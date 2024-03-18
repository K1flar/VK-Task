package filmservice

import (
	"film_library/internal/config"
	"film_library/internal/domains"
	"film_library/pkg/pagination"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrInvalidName        = fmt.Errorf("invalid film name")
	ErrInvalidDescription = fmt.Errorf("invalid film description")
	ErrInvalidRating      = fmt.Errorf("invalid film rating")
)

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

type ActorService interface {
	AddActorsToFilm(filmID uint32, actors []uint32) error
}

type FilmService struct {
	repo         FilmRepo
	actorService ActorService
	log          *slog.Logger
	cfg          *config.Config
}

func New(repo FilmRepo, actorService ActorService, log *slog.Logger, cfg *config.Config) *FilmService {
	return &FilmService{
		repo:         repo,
		actorService: actorService,
		log:          log,
		cfg:          cfg,
	}
}

func (s *FilmService) CreateFilm(film domains.Film, actorsID []uint32) (uint32, error) {
	fn := "filmService.CreateFilm"

	err := s.validateFilm(film)

	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return 0, err
	}

	filmID, err := s.repo.AddFilm(film)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	err = s.actorService.AddActorsToFilm(filmID, actorsID)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return filmID, nil
}

func (s *FilmService) UpdateFilmName(id uint32, name string) error {
	fn := "filmService.UpdateFilmName"

	minNameLen, maxNameLen := s.cfg.FilmValidations.MinNameLen, s.cfg.FilmValidations.MaxNameLen
	if len(name) < minNameLen || len(name) > maxNameLen {
		s.log.Error(fmt.Sprintf("%s: %s", fn, ErrInvalidName.Error()))
		return fmt.Errorf("%s: %w", fn, ErrInvalidName)
	}

	err := s.repo.UpdateFilmName(id, name)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *FilmService) UpdateFilmDescription(id uint32, descrtion string) error {
	fn := "filmService.UpdateFilmName"

	minDescriptionLen, maxDescriptionLen := s.cfg.FilmValidations.MinDescriptionLen, s.cfg.FilmValidations.MaxDescriptionLen
	if len(descrtion) < minDescriptionLen || len(descrtion) > maxDescriptionLen {
		s.log.Error(fmt.Sprintf("%s: %s", fn, ErrInvalidDescription.Error()))
		return fmt.Errorf("%s: %w", fn, ErrInvalidDescription)
	}

	err := s.repo.UpdateFilmDescription(id, descrtion)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *FilmService) UpdateFilmReleaseDate(id uint32, releaseDate time.Time) error {
	fn := "filmService.UpdateFilmReleaseDate"

	err := s.repo.UpdateFilmReleaseDate(id, releaseDate)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *FilmService) UpdateFilmRating(id uint32, rating int) error {
	fn := "filmService.UpdateFilmRating"

	minRating, maxRating := s.cfg.FilmValidations.MinRating, s.cfg.FilmValidations.MaxRating
	fmt.Println(rating, minRating, maxRating)
	if rating < minRating || rating > maxRating {
		s.log.Error(fmt.Sprintf("%s: %s", fn, ErrInvalidRating.Error()))
		return fmt.Errorf("%s: %w", fn, ErrInvalidRating)
	}

	err := s.repo.UpdateFilmRating(id, rating)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *FilmService) UpdateFilm(id uint32, film domains.Film) error {
	fn := "filmService.UpdateFilm"

	err := s.validateFilm(film)

	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return err
	}

	err = s.repo.UpdateFilm(id, film)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *FilmService) DeleteFilm(id uint32) error {
	fn := "filmService.DeleteFilm"
	err := s.repo.DeleteFilm(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *FilmService) GetFilms(filter *pagination.FilmFilter) ([]*domains.Film, error) {
	fn := "filmService.GetFilms"

	filter.Validate()

	films, err := s.repo.GetFilms(filter)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return films, nil
}
