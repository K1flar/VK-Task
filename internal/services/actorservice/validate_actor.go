package actorservice

import (
	"film_library/internal/domains"
	"film_library/pkg/validation"
)

func (s *ActorService) validateActor(actor domains.Actor) error {
	err := validation.NewValidator[domains.Actor](actor).
		Must(
			func(a domains.Actor) bool { return len(actor.FullName) > 0 },
			ErrInvalidFullName.Error()).
		Must(
			func(a domains.Actor) bool { return actor.Gender.IsValid() },
			ErrInvalidGender.Error()).
		Validate()

	return err
}
