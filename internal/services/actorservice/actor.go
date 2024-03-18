package actorservice

import (
	"film_library/internal/domains"
	"film_library/pkg/pagination"
	"film_library/pkg/validation"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrInvalidFullName = fmt.Errorf("full name must be at least 1 letter long")
	ErrInvalidGender   = fmt.Errorf("gender must be male or female")
)

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

type ActorService struct {
	repo ActorRepo
	log  *slog.Logger
}

func New(repo ActorRepo, log *slog.Logger) *ActorService {
	return &ActorService{
		repo: repo,
		log:  log,
	}
}

func (s *ActorService) CreateActor(actor domains.Actor) error {
	fn := "actorService.CreateActor"

	err := s.validateActor(actor)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return err
	}

	err = s.repo.AddActor(actor)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) AddActorsToFilm(filmID uint32, actorsID []uint32) error {
	fn := "actorService.AddActorsToFilm"

	err := s.repo.AddActorsToFilm(filmID, actorsID)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) UpdateActorFullName(id uint32, fullName string) error {
	fn := "actorService.UpdateActorFullName"

	if len(fullName) == 0 {
		return fmt.Errorf("%s: %w", fn, ErrInvalidFullName)
	}

	err := s.repo.UpdateActorFullName(id, fullName)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) UpdateActorGender(id uint32, gender domains.Gender) error {
	fn := "actorService.UpdateActorGender"

	if !gender.IsValid() {
		s.log.Error(fmt.Sprintf("%s: %s: %s", fn, ErrInvalidGender.Error(), gender))
		return fmt.Errorf("%s: %w", fn, ErrInvalidGender)
	}

	err := s.repo.UpdateActorGender(id, string(gender))
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) UpdateActorBirthday(id uint32, birthday time.Time) error {
	fn := "actorService.UpdateActorBirthday"
	err := s.repo.UpdateActorBirthday(id, birthday)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) UpdateActor(id uint32, actor domains.Actor) error {
	fn := "actorService.UpdateActor"

	err := validation.NewValidator[domains.Actor](actor).
		Must(
			func(a domains.Actor) bool { return len(actor.FullName) > 0 },
			ErrInvalidFullName.Error()).
		Must(
			func(a domains.Actor) bool { return actor.Gender.IsValid() },
			ErrInvalidGender.Error()).
		Validate()

	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	err = s.repo.UpdateActor(id, actor)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) DeleteActor(id uint32) error {
	fn := "actorService.DeleteActor"
	err := s.repo.DeleteActor(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) DeleteActorFromFilm(actorID uint32, filmID uint32) error {
	fn := "actorService.DeleteActor"
	err := s.repo.DeleteActorFromFilm(actorID, filmID)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *ActorService) GetActorsWithFilms(filter *pagination.ActorsFilter) ([]*domains.ActorWithFilms, error) {
	fn := "actorService.GetActorsWithFilms"

	filter.Pagination.ValidatePagination()

	actorWithFilms, err := s.repo.GetActorsWithFilms(filter)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return actorWithFilms, nil
}
