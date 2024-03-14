package actorrepo

import (
	"database/sql"
	"film_library/internal/domains"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var (
	ErrInvalidGender = fmt.Errorf("invalid actor gender")
	ErrNotFound      = fmt.Errorf("actor not found")
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{
		db: db,
	}
}

func (r *ActorRepository) AddActor(actor domains.Actor) error {
	fn := "actorRepository.AddActor"

	stmt := `
		INSERT INTO actors(full_name, gender, birthday)
		VALUE ($1, $2, $3);
	`

	_, err := r.db.Exec(stmt, actor.FullName, actor.Gender, actor.Birthday)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidGender)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) updateFieldByID(id uint32, field string, value any) error {
	stmt := fmt.Sprintf(`
		UPDATE actors
		SET %s=$1
		WHERE id=$2
	`, field)

	res, err := r.db.Exec(stmt, value, id)
	if err != nil {
		return err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *ActorRepository) UpdateActorFullName(id uint32, fullName string) error {
	fn := "actorRepository.UpdateActorFullName"

	if err := r.updateFieldByID(id, "full_name", fullName); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) UpdateActorGender(id uint32, gender string) error {
	fn := "actorRepository.UpdateActorGender"

	if err := r.updateFieldByID(id, "gender", gender); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidGender)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) UpdateActorBirthday(id uint32, birthday time.Time) error {
	fn := "actorRepository.UpdateActorBirthday"

	if err := r.updateFieldByID(id, "birthday", birthday); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
