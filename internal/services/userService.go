package services

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"workout-tracker/internal/models"
	"workout-tracker/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrRequiredFields = errors.New("required fields are empty")
var ErrInvalidCredentials = errors.New("invalid credentials")

var secretKey = []byte(os.Getenv("JWT_KEY"))

type UserService struct {
	Repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (service *UserService) validateCredentials(name, password string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(password) == "" {
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

func (service *UserService) createToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
			"sub":      strconv.FormatInt(user.Id, 10),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (service *UserService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		if len(secretKey) == 0 {
			return nil, fmt.Errorf("jwt secret key is not set")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (service *UserService) Login(name, password string) (string, error) { // returns jwt as string
	if err := service.validateCredentials(name, password); err != nil {
		return "", err
	}
	user, err := service.Repo.GetUserByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	if err := service.comparePasswords(password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}
	token, err := service.createToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *UserService) Register(name, password string) (string, error) {
	if err := service.validateCredentials(name, password); err != nil {
		return "", err
	}
	passwordHash, err := service.hashPassword(password)
	if err != nil {
		return "", err
	}
	id, err := service.Repo.InsertUser(name, passwordHash)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil

}
