package services

import (
	"database/sql"
	"errors"
	"strings"
	"workout-tracker/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

var ErrRequiredFields = errors.New("required fields are empty")
var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	Repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (service *UserService) validateCredentials(name, password string) error {
	if strings.TrimSpace(name) == "" && strings.TrimSpace(password) == "" {
		return ErrRequiredFields
	}
	return nil
}

func (service *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (service *UserService) comparePasswords(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash)); err != nil {
		return err
	}
	return nil
}

func (service *UserService) Login(name, password string) (string, error) { // returns jwt as string
	if err := service.validateCredentials(name, password); err != nil {
		return "", err
	}
	user, err := service.Repo.GetUserByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrRequiredFields
		}
		return "", err
	}
	if err := service.comparePasswords(password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}
	// return token
	return user.Username, nil // placeholder
}
