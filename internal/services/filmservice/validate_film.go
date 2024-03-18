package filmservice

import (
	"film_library/internal/domains"
	"film_library/pkg/validation"
)

func (s *FilmService) validateFilm(film domains.Film) error {
	minNameLen, maxNameLen := s.cfg.FilmValidations.MinNameLen, s.cfg.FilmValidations.MaxNameLen
	minDescriptionLen, maxDescriptionLen := s.cfg.FilmValidations.MinDescriptionLen, s.cfg.FilmValidations.MaxDescriptionLen
	minRating, maxRating := s.cfg.FilmValidations.MinRating, s.cfg.FilmValidations.MaxRating

	err := validation.NewValidator[domains.Film](film).
		Between(
			func(f domains.Film) int { return len(f.Name) },
			minNameLen, maxNameLen,
			ErrInvalidName.Error()).
		Between(
			func(f domains.Film) int { return len(f.Description) },
			minDescriptionLen, maxDescriptionLen,
			ErrInvalidDescription.Error()).
		Between(
			func(f domains.Film) int { return f.Rating },
			minRating, maxRating,
			ErrInvalidRating.Error()).
		Validate()

	return err
}
