package services

import (
	"database/sql"
	"errors"
	"strings"
	"workout-tracker/internal/repo"
)

var ErrRequiredFields = errors.New("required fields are empty")

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
	// check if passwords match
	// return token

}
