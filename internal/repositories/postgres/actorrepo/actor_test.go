package actorrepo

import (
	"errors"
	"film_library/internal/domains"
	"film_library/pkg/pagination"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestActorRepoAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(actor domains.Actor)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name  string
		actor domains.Actor
		mock  mockBehavior
		err   error
	}{
		{
			name:  "Correct",
			actor: domains.Actor{FullName: "Robert Oppenheimer", Gender: "male", Birthday: domains.Time(time.Now())},
			mock: func(actor domains.Actor) {
				mock.ExpectExec("INSERT INTO actors").
					WithArgs(actor.FullName, actor.Gender, actor.Birthday).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:  "Invalid gender",
			actor: domains.Actor{FullName: "denis", Gender: "male2", Birthday: domains.Time(time.Now())},
			mock: func(actor domains.Actor) {
				mock.ExpectExec("INSERT INTO actors").
					WithArgs(actor.FullName, actor.Gender, actor.Birthday).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23514")})
			},
			err: ErrInvalidGender,
		},
		{
			name:  "Unknown error",
			actor: domains.Actor{FullName: "123", Gender: "123", Birthday: domains.Time(time.Now())},
			mock: func(actor domains.Actor) {
				mock.ExpectExec("INSERT INTO actors").
					WithArgs(actor.FullName, actor.Gender, actor.Birthday).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.actor)

			err := repo.AddActor(tc.actor)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActorRepoUpdateFullName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(id uint32, fullName string)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name     string
		id       uint32
		fullName string
		mock     mockBehavior
		err      error
	}{
		{
			name:     "Correct",
			id:       1,
			fullName: "Robert",
			mock: func(id uint32, fullName string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(fullName, id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:     "Not found",
			id:       1,
			fullName: "denis",
			mock: func(id uint32, fullName string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(fullName, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			err: ErrNotFound,
		},
		{
			name:     "Unknown error",
			id:       1,
			fullName: "aboba",
			mock: func(id uint32, fullName string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(fullName, id).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id, tc.fullName)

			err := repo.UpdateActorFullName(tc.id, tc.fullName)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActorRepoUpdateGender(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(id uint32, gender string)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name   string
		id     uint32
		gender string
		mock   mockBehavior
		err    error
	}{
		{
			name:   "Correct",
			id:     1,
			gender: "female",
			mock: func(id uint32, gender string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(gender, id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:   "Not found",
			id:     156,
			gender: "male",
			mock: func(id uint32, gender string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(gender, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			err: ErrNotFound,
		},
		{
			name:   "Invalid gender",
			id:     1,
			gender: "aboba",
			mock: func(id uint32, fullName string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(fullName, id).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23514")})
			},
			err: ErrInvalidGender,
		},
		{
			name:   "Unknown error",
			id:     2,
			gender: "male",
			mock: func(id uint32, fullName string) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(fullName, id).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id, tc.gender)

			err := repo.UpdateActorGender(tc.id, tc.gender)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActorRepoUpdateBirthday(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(id uint32, birthday time.Time)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name     string
		id       uint32
		birthday time.Time
		mock     mockBehavior
		err      error
	}{
		{
			name:     "Correct",
			id:       1,
			birthday: time.Now(),
			mock: func(id uint32, birthday time.Time) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(birthday, id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:     "Not found",
			id:       156,
			birthday: time.Now(),
			mock: func(id uint32, birthday time.Time) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(birthday, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			err: ErrNotFound,
		},
		{
			name:     "Unknown error",
			id:       2,
			birthday: time.Now(),
			mock: func(id uint32, birthday time.Time) {
				mock.ExpectExec("UPDATE actors").
					WithArgs(birthday, id).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id, tc.birthday)

			err := repo.UpdateActorBirthday(tc.id, tc.birthday)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActorRepoDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(id uint32)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name string
		id   uint32
		mock mockBehavior
		err  error
	}{
		{
			name: "Correct",
			id:   1,
			mock: func(id uint32) {
				mock.ExpectExec("DELETE FROM actors").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Not found",
			id:   156,
			mock: func(id uint32) {
				mock.ExpectExec("DELETE FROM actors").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			err: ErrNotFound,
		},
		{
			name: "Not found",
			id:   156,
			mock: func(id uint32) {
				mock.ExpectExec("DELETE FROM actors").
					WithArgs(id).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id)

			err := repo.DeleteActor(tc.id)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActorRepoAddActors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(filmID uint32, actorsID []uint32)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name     string
		filmID   uint32
		actorsID []uint32
		mock     mockBehavior
		err      error
	}{
		{
			name:     "Correct",
			filmID:   1,
			actorsID: []uint32{1, 2, 3, 4},
			mock: func(filmID uint32, actorsID []uint32) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO film_actor(actor_id, film_id) VALUES (1, 1),(2, 1),(3, 1),(4, 1);")).
					WillReturnResult(sqlmock.NewResult(0, 3))
			},
		},
		{
			name:     "No actors",
			filmID:   1,
			actorsID: []uint32{},
			mock:     func(filmID uint32, actorsID []uint32) {},
		},
		{
			name:     "Not unique actors",
			filmID:   1,
			actorsID: []uint32{2, 2},
			mock: func(filmID uint32, actorsID []uint32) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO film_actor(actor_id, film_id) VALUES (2, 1),(2, 1);")).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23505"), Constraint: "film_actor_pkey"})
			},
			err: ErrUniqueActors,
		},
		{
			name:     "Actor not found",
			filmID:   1,
			actorsID: []uint32{1024},
			mock: func(filmID uint32, actorsID []uint32) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO film_actor(actor_id, film_id) VALUES (1024, 1);")).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23503"), Constraint: "film_actor_actor_id_fkey"})
			},
			err: ErrNotFound,
		},
		{
			name:     "Film not found",
			filmID:   1,
			actorsID: []uint32{1024},
			mock: func(filmID uint32, actorsID []uint32) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO film_actor(actor_id, film_id) VALUES (1024, 1);")).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23503"), Constraint: "film_actor_film_id_fkey"})
			},
			err: ErrNotFound,
		},
		{
			name:     "Unknown error",
			filmID:   2,
			actorsID: []uint32{2},
			mock: func(filmID uint32, actorsID []uint32) {
				mock.ExpectExec("INSERT INTO film_actor").
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.filmID, tc.actorsID)

			err := repo.AddActorsToFilm(tc.filmID, tc.actorsID)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActorRepoGetActorsWithFilms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewActorRepository(db)

	type mockBehavior func(filter *pagination.ActorsFilter)

	// customError := fmt.Errorf("some error")
	tests := []struct {
		name            string
		filter          *pagination.ActorsFilter
		mock            mockBehavior
		actorsWithFilms []*domains.ActorWithFilms
		err             error
	}{
		{
			name: "Correct",
			filter: &pagination.ActorsFilter{
				Pagination:       pagination.New(1, 10),
				FullNameContains: "Rob",
			},
			mock: func(filter *pagination.ActorsFilter) {
				rows := sqlmock.NewRows([]string{"id", "full_name", "gender", "birthday", "id", "name", "description", "release_date", "rating"}).
					AddRow(1, "Roby", "male", time.Now(), 1, "Oppenheimer", "", time.Now(), 10).
					AddRow(1, "Roby", "male", time.Now(), 10, "Abobaheimer", "", time.Now(), 9).
					AddRow(2, "Aboba", "female", time.Now(), 10, "Abobaheimer", "", time.Now(), 9)
				mock.ExpectQuery(`SELECT a.id, a.full_name, a.gender, a.birthday, 
					f.id, f.name, f.description, f.release_date, f.rating FROM actors AS a`).
					WithArgs(strings.ToLower(filter.FullNameContains)).
					WillReturnRows(rows)
			},
			actorsWithFilms: []*domains.ActorWithFilms{
				{
					Actor: domains.Actor{ID: 1},
					Films: []*domains.Film{{ID: 1}, {ID: 10}},
				},
				{
					Actor: domains.Actor{ID: 2},
					Films: []*domains.Film{{ID: 10}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.filter)

			got, err := repo.GetActorsWithFilms(tc.filter)

			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("expected: %s\ngot: %s", tc.err, err)
				}
			} else {
				if len(got) != len(tc.actorsWithFilms) {
					t.Errorf("expected: %#v\ngot: %#v", len(tc.actorsWithFilms), len(got))
				}
				for i := 0; i < len(got); i++ {
					if got[i].ID != tc.actorsWithFilms[i].ID {
						t.Errorf("expected: %#v\ngot: %#v", tc.actorsWithFilms, got)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
