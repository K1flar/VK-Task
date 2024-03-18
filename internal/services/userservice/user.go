package userservice

import (
	"errors"
	"film_library/internal/config"
	"film_library/internal/domains"
	"film_library/internal/repositories/postgres/userrepo"
	"film_library/pkg/validation"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRole        = fmt.Errorf("invalid role")
	ErrInvalidLoginLen    = fmt.Errorf("invalid login length")
	ErrInvalidPasswordLen = fmt.Errorf("invalid password length")
	ErrNotFound           = fmt.Errorf("user not found")
	ErrInvalidPassword    = fmt.Errorf("invalid password")
)

type UserRepo interface {
	AddUser(user domains.User) error
	GetUserByLoign(login string) (*domains.User, error)
}

type UserService struct {
	repo UserRepo
	log  *slog.Logger
	cfg  *config.Config
}

func New(repo UserRepo, log *slog.Logger, cfg *config.Config) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *UserService) CreateUser(user domains.User) (string, error) {
	fn := "userService.CreateUser"

	minLoginLen, maxLoginLen := s.cfg.Identity.MinLoginLen, s.cfg.Identity.MaxLoginLen
	minPasswordLen, maxPasswordLen := s.cfg.Identity.MinPasswordLen, s.cfg.Identity.MaxPasswordLen

	err := validation.NewValidator[domains.User](user).
		Between(
			func(u domains.User) int { return len(u.Login) },
			minLoginLen, maxLoginLen,
			ErrInvalidLoginLen.Error()).
		Between(
			func(u domains.User) int { return len(u.Password) },
			minPasswordLen, maxPasswordLen,
			ErrInvalidPasswordLen.Error()).
		Must(
			func(u domains.User) bool { return user.Role.IsValidRole() },
			ErrInvalidRole.Error()).
		Validate()

	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return "", err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: error occurred generating hash password: %s", fn, ErrInvalidRole.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	user.Password = string(hashPassword)

	err = s.repo.AddUser(user)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return s.generateToken(&user, s.cfg.Server.Secret)
}

func (s *UserService) GetUserByLogin(login string) (*domains.User, error) {
	fn := "userService.GetUserByLogin"

	u, err := s.repo.GetUserByLoign(login)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return u, nil
}

func (s *UserService) Login(login, password string) (string, error) {
	fn := "userService.Login"

	user, err := s.repo.GetUserByLoign(login)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		if errors.Is(err, userrepo.ErrNotFound) {
			return "", fmt.Errorf("%s: %w", fn, ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Warn(fmt.Sprintf("%s: %s", fn, err.Error()))
		return "", fmt.Errorf("%s: %w", fn, ErrInvalidPassword)
	}

	return s.generateToken(user, s.cfg.Server.Secret)
}

func (s *UserService) generateToken(user *domains.User, secret string) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"login": user.Login,
		"role":  user.Role,
	})

	return jwt.SignedString([]byte(secret))
}
