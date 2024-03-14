package actorrepo

import (
	"errors"
	"film_library/internal/domains"
	"fmt"
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
			actor: domains.Actor{FullName: "Robert Oppenheimer", Gender: "male", Birthday: time.Now()},
			mock: func(actor domains.Actor) {
				mock.ExpectExec("INSERT INTO actors").
					WithArgs(actor.FullName, actor.Gender, actor.Birthday).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:  "Invalid gender",
			actor: domains.Actor{FullName: "denis", Gender: "male2", Birthday: time.Now()},
			mock: func(actor domains.Actor) {
				mock.ExpectExec("INSERT INTO actors").
					WithArgs(actor.FullName, actor.Gender, actor.Birthday).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23514")})
			},
			err: ErrInvalidGender,
		},
		{
			name:  "Unknown error",
			actor: domains.Actor{FullName: "123", Gender: "123", Birthday: time.Now()},
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
