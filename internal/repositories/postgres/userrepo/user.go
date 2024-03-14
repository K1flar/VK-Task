package userrepo

import (
	"database/sql"
	"film_library/internal/domains"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrInvalidRole  = fmt.Errorf("invalid user role")
	ErrAlredyExists = fmt.Errorf("user alredy exists")
	ErrNotFound     = fmt.Errorf("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddUser(user domains.User) error {
	fn := "userRepository.AddUser"

	stmt := `
		INSERT INTO users(login, password, role)
		VALUES ($1, $2, $3);
	`

	_, err := r.db.Exec(stmt, user.Login, user.Password, user.Role)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch {
			case err.Code == pq.ErrorCode("23514"):
				return fmt.Errorf("%s: %w", fn, ErrInvalidRole)
			case err.Code == pq.ErrorCode("23505"):
				return fmt.Errorf("%s: %w", fn, ErrAlredyExists)
			}
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *UserRepository) GetUserByLoign(login string) (*domains.User, error) {
	fn := "userRepository.GetUserByLoign"

	stmt := `
		SELECT id, login, password, role
		FROM users
		WHERE login=$1
	`

	user := &domains.User{}
	row := r.db.QueryRow(stmt, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}
