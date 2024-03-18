package userrepo

import (
	"errors"
	"film_library/internal/domains"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestUserRepoAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type mockBehavior func(user domains.User)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name string
		user domains.User
		mock mockBehavior
		err  error
	}{
		{
			name: "Correct",
			user: domains.User{Login: "denis", Password: "denis", Role: "admin"},
			mock: func(user domains.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Login, user.Password, user.Role).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Already exists",
			user: domains.User{Login: "denis", Password: "denis", Role: "admin"},
			mock: func(user domains.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Login, user.Password, user.Role).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23505")})
			},
			err: ErrAlreadyExists,
		},
		{
			name: "Invalid role",
			user: domains.User{Login: "aboba", Password: "admin", Role: "abobavich"},
			mock: func(user domains.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Login, user.Password, user.Role).
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23514")})
			},
			err: ErrInvalidRole,
		},
		{
			name: "Unknown error",
			user: domains.User{Login: "123", Password: "123", Role: "viewer"},
			mock: func(user domains.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Login, user.Password, user.Role).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.user)

			err := repo.AddUser(tc.user)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %s\ngot: %s", tc.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepoGetByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type mockBehavior func(login string)

	customError := fmt.Errorf("some error")
	tests := []struct {
		name  string
		login string
		mock  mockBehavior
		user  *domains.User
		err   error
	}{
		{
			name:  "Correct",
			login: "denis",
			mock: func(login string) {
				rows := mock.NewRows([]string{"id", "login", "password", "role"}).AddRow(1, "denis", "denis", "admin")
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(login).
					WillReturnRows(rows)
			},
			user: &domains.User{ID: 1},
		},
		{
			name:  "Not found",
			login: "denis",
			mock: func(login string) {
				rows := mock.NewRows([]string{"id", "login", "password", "role"})
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(login).
					WillReturnRows(rows)
			},
			err: ErrNotFound,
		},
		{
			name:  "Unknown error",
			login: "denis",
			mock: func(login string) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(login).
					WillReturnError(customError)
			},
			err: customError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.login)

			got, err := repo.GetUserByLoign(tc.login)

			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("expected: %s\ngot: %s", tc.err, err)
				}
			} else {
				if got == nil || got.ID != tc.user.ID {
					t.Errorf("expected: %#v\ngot: %#v", tc.user, got)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
