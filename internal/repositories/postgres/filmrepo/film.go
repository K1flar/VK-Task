package filmrepo

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
	ErrNotFound          = fmt.Errorf("film not found")
	ErrInvalidNameLength = fmt.Errorf("invalid film name length")
	ErrInvalidRating     = fmt.Errorf("invalid film rating")
	ErrAlreadyExists     = fmt.Errorf("film already exists")
)

type FilmRepository struct {
	db *sql.DB
}

func NewFilmRepository(db *sql.DB) *FilmRepository {
	return &FilmRepository{
		db: db,
	}
}

func (r *FilmRepository) AddFilm(film domains.Film) (uint32, error) {
	fn := "filmRepository.AddFilm"

	stmt := `
		INSERT INTO films(name, description, release_date, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var filmID int
	row := r.db.QueryRow(stmt, film.Name, film.Description, time.Time(film.ReleaseDate), film.Rating)
	err := row.Scan(&filmID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Constraint {
			case "films_name_key":
				return 0, fmt.Errorf("%s: %w", fn, ErrAlreadyExists)
			case "films_name_check":
				return 0, fmt.Errorf("%s: %w", fn, ErrInvalidNameLength)
			case "films_rating_check":
				return 0, fmt.Errorf("%s: %w", fn, ErrInvalidRating)
			}
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return uint32(filmID), nil
}

func (r *FilmRepository) updateField(id uint32, field string, value any) error {
	stmt := fmt.Sprintf(`
		UPDATE films
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

func (r *FilmRepository) UpdateFilmName(id uint32, name string) error {
	fn := "filmRepository.UpdateFilmName"
	if err := r.updateField(id, "name", name); err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Constraint {
			case "films_name_key":
				return fmt.Errorf("%s: %w", fn, ErrAlreadyExists)
			case "films_name_check":
				return fmt.Errorf("%s: %w", fn, ErrInvalidNameLength)
			}
		}
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *FilmRepository) UpdateFilmDescription(id uint32, description string) error {
	fn := "filmRepository.UpdateFilmDescription"
	if err := r.updateField(id, "description", description); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *FilmRepository) UpdateFilmReleaseDate(id uint32, releaseDate time.Time) error {
	fn := "filmRepository.UpdateFilmReleaseDate"
	if err := r.updateField(id, "release_date", releaseDate); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *FilmRepository) UpdateFilmRating(id uint32, rating int) error {
	fn := "filmRepository.UpdateFilmReleaseDate"
	if err := r.updateField(id, "rating", rating); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pq.ErrorCode("23514") {
			return fmt.Errorf("%s: %w", fn, ErrInvalidRating)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (r *FilmRepository) UpdateFilm(id uint32, film domains.Film) error {
	fn := "actorRepository.UpdateFilm"

	stmt := `
		UPDATE films
		SET (name, description, release_date, rating) = ($1, $2, $3, $4)
		WHERE id=$5;
	`
	res, err := r.db.Exec(stmt, film.Name, film.Description, time.Time(film.ReleaseDate), film.Rating, id)
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

func (r *FilmRepository) DeleteFilm(id uint32) error {
	fn := "filmRepository.DeleteFilm"

	stmt := `
		DELETE FROM films
		WHERE id=$1
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

func (r *FilmRepository) GetFilms(filter *pagination.FilmFilter) ([]*domains.Film, error) {
	fn := "filmRepository.GetFilms"

	query := selectbuilder.New("SELECT DISTINCT f.id, f.name, f.description, f.release_date, f.rating FROM films AS f")
	if filter.ActorNameContains != "" {
		query.Join("film_actor AS fa ON f.id=fa.film_id").
			Join("actors AS a ON a.id=fa.actor_id").
			Where("LOWER(a.full_name) LIKE '%s'", "%"+strings.ToLower(filter.ActorNameContains)+"%")
	}

	q := query.Where("LOWER(f.name) LIKE '%s'", "%"+strings.ToLower(filter.NameContains)+"%").
		OrderBy(filter.OrderBy, filter.Direction).
		AddPagination(filter.Pagination).
		Build()

	res, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	films := []*domains.Film{}
	for res.Next() {
		film := &domains.Film{}
		err = res.Scan(&film.ID, &film.Name, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		films = append(films, film)
	}

	return films, nil
}
