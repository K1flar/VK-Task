package filmrepo

import (
	"errors"
	"film_library/internal/domains"
	"film_library/pkg/pagination"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestFilmRepoAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewFilmRepository(db)

	type mockBehavior func(film domains.Film)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name string
		film domains.Film
		mock mockBehavior
		id   uint32
		err  error
	}{
		{
			name: "Correct",
			film: domains.Film{Name: "Oppenheimer", ReleaseDate: domains.Time(time.Now()), Rating: 10},
			mock: func(film domains.Film) {
				mock.ExpectExec("INSERT INTO films").
					WithArgs(film.Name, film.Description, film.ReleaseDate, film.Rating).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: 1,
		},
		{
			name: "Already exists",
			film: domains.Film{Name: "Oppenheimer", ReleaseDate: domains.Time(time.Now()), Rating: 10},
			mock: func(film domains.Film) {
				mock.ExpectExec("INSERT INTO films").
					WithArgs(film.Name, film.Description, film.ReleaseDate, film.Rating).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23505"), Constraint: "films_name_key"})
			},
			err: ErrAlreadyExists,
		},
		{
			name: "Invalid name",
			film: domains.Film{ReleaseDate: domains.Time(time.Now()), Rating: 10},
			mock: func(film domains.Film) {
				mock.ExpectExec("INSERT INTO films").
					WithArgs(film.Name, film.Description, film.ReleaseDate, film.Rating).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23514"), Constraint: "films_name_check"})
			},
			err: ErrInvalidNameLength,
		},
		{
			name: "Invalid rating",
			film: domains.Film{Name: "Barby", ReleaseDate: domains.Time(time.Now()), Rating: -1},
			mock: func(film domains.Film) {
				mock.ExpectExec("INSERT INTO films").
					WithArgs(film.Name, film.Description, film.ReleaseDate, film.Rating).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23505"), Constraint: "films_rating_check"})
			},
			err: ErrInvalidRating,
		},
		{
			name: "Unknown error",
			film: domains.Film{Name: "Aboba", ReleaseDate: domains.Time(time.Now()), Rating: 5},
			mock: func(film domains.Film) {
				mock.ExpectExec("INSERT INTO films").
					WithArgs(film.Name, film.Description, film.ReleaseDate, film.Rating).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.film)

			got, err := repo.AddFilm(tc.film)

			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("expected: %s\ngot: %s", tc.err, err)
				}
			} else {
				if got == 0 || got != tc.id {
					t.Errorf("expected: %#v\ngot: %#v", tc.id, got)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFilmRepoGetFilms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewFilmRepository(db)

	type mockBehavior func(filter *pagination.FilmFilter)

	// customError := fmt.Errorf("some error")
	tests := []struct {
		name   string
		filter *pagination.FilmFilter
		mock   mockBehavior
		films  []*domains.Film
		err    error
	}{
		{
			name: "Correct",
			filter: &pagination.FilmFilter{
				Pagination:        pagination.New(1, 10),
				NameContains:      "oppen",
				ActorNameContains: "rob",
			},
			mock: func(filter *pagination.FilmFilter) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).
					AddRow(1, "Oppenheimer", "", time.Now(), 10)
				mock.ExpectQuery("SELECT DISTINCT f.id, f.name, f.description, f.release_date, f.rating FROM films AS f").
					WithArgs(strings.ToLower(filter.ActorNameContains), strings.ToLower(filter.NameContains)).
					WillReturnRows(rows)
			},
			films: []*domains.Film{{ID: 1, Name: "Oppenheimer", ReleaseDate: domains.Time(time.Now()), Rating: 10}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.filter)

			got, err := repo.GetFilms(tc.filter)

			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("expected: %s\ngot: %s", tc.err, err)
				}
			} else {
				if len(got) != len(tc.films) {
					t.Errorf("expected: %#v\ngot: %#v", len(tc.films), len(got))
				}
				for i := 0; i < len(got); i++ {
					if got[i].ID != tc.films[i].ID {
						t.Errorf("expected: %#v\ngot: %#v", tc.films, got)
					}
				}
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
