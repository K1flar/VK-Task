package actorrepo

import (
	"database/sql"
	"film_library/internal/domains"
	"film_library/pkg/pagination"
	selectbuilder "film_library/pkg/sqltools/select_builder"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	ErrInvalidGender = fmt.Errorf("invalid actor gender")
	ErrNotFound      = fmt.Errorf("actor not found")
	ErrUniqueActors  = fmt.Errorf("actors must be unique")
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
		VALUES ($1, $2, $3);
	`

	_, err := r.db.Exec(stmt, actor.FullName, actor.Gender, time.Time(actor.Birthday))
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidGender)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) updateField(id uint32, field string, value any) error {
	stmt := fmt.Sprintf(`
		UPDATE actors
		SET %s=$1
		WHERE id=$2;
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

	if err := r.updateField(id, "full_name", fullName); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) UpdateActorGender(id uint32, gender string) error {
	fn := "actorRepository.UpdateActorGender"

	if err := r.updateField(id, "gender", gender); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidGender)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) UpdateActorBirthday(id uint32, birthday time.Time) error {
	fn := "actorRepository.UpdateActorBirthday"

	if err := r.updateField(id, "birthday", birthday); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *ActorRepository) UpdateActor(id uint32, actor domains.Actor) error {
	fn := "actorRepository.UpdateActor"

	stmt := `
		UPDATE actors
		SET (full_name, gender, birthday)=($1, $2, $3)
		WHERE id=$4;
	`

	res, err := r.db.Exec(stmt, actor.FullName, actor.Gender, time.Time(actor.Birthday), id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		return fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	return nil
}

func (r *ActorRepository) DeleteActor(id uint32) error {
	fn := "actorRepository.DeleteActor"

	stmt := `
		DELETE FROM actors
		WHERE id=$1;
	`

	res, err := r.db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		return fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	return nil
}

func (r *ActorRepository) DeleteActorFromFilm(actorID uint32, filmID uint32) error {
	fn := "actorRepository.DeleteActorFromFilm"

	stmt := `
		DELETE FROM film_actor
		WHERE film_id=$1 AND actor_id=$2;
	`

	res, err := r.db.Exec(stmt, filmID, actorID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAff == 0 {
		return fmt.Errorf("%s: %w", fn, ErrNotFound)
	}

	return nil
}

func (r *ActorRepository) GetActorsWithFilms(filter *pagination.ActorsFilter) ([]*domains.ActorWithFilms, error) {
	fn := "actorRepository.GetActorsWithFilms"
	query := selectbuilder.
		New(`SELECT a.id, a.full_name, a.gender, a.birthday, 
			f.id, f.name, f.description, f.release_date, f.rating FROM actors AS a`).
		Join("film_actor AS fa ON a.id=fa.actor_id").
		Join("films AS f ON f.id=fa.film_id").
		Where("LOWER(a.full_name) LIKE $1").
		AddPagination(filter.Pagination).
		Build()

	res, err := r.db.Query(query, "%"+strings.ToLower(filter.FullNameContains)+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	actorsWithFilms := []*domains.ActorWithFilms{}
	indexesOfActors := map[uint32]int{}

	for res.Next() {
		actor := &domains.Actor{}
		film := &domains.Film{}
		err := res.Scan(&actor.ID, &actor.FullName, &actor.Gender, &actor.Birthday,
			&film.ID, &film.Name, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		if _, ok := indexesOfActors[actor.ID]; !ok {
			actorsWithFilms = append(actorsWithFilms, &domains.ActorWithFilms{Actor: *actor, Films: []*domains.Film{}})
			indexesOfActors[actor.ID] = len(actorsWithFilms) - 1
		}
		films := &actorsWithFilms[indexesOfActors[actor.ID]].Films
		*films = append(*films, film)
	}

	return actorsWithFilms, nil
}

func (r *ActorRepository) AddActorsToFilm(filmID uint32, actorsID []uint32) error {
	fn := "actorRepository.AddActorsToFilm"

	if len(actorsID) == 0 {
		return nil
	}

	var rowsBuilder strings.Builder
	for _, id := range actorsID {
		_, err := rowsBuilder.WriteString(fmt.Sprintf("(%d, %d),", id, filmID))
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}
	rows := rowsBuilder.String()
	rows = rows[:len(rows)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO film_actor(actor_id, film_id) 
		VALUES %s;
	`, rows)

	_, err := r.db.Exec(stmt)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Constraint {
			case "film_actor_pkey":
				return fmt.Errorf("%s: %w", fn, ErrUniqueActors)
			case "film_actor_actor_id_fkey":
				return fmt.Errorf("%s: %w", fn, ErrNotFound)
			case "film_actor_film_id_fkey":
				return fmt.Errorf("%s: %w", fn, ErrNotFound)
			}
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
