package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
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

var ErrMissingRefreshToken = errors.New("missing refresh token")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type UserService struct {
	Repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (service *UserService) getSecretKey() ([]byte, error) {
	key := strings.TrimSpace(os.Getenv("JWT_KEY"))
	if key == "" {
		return nil, fmt.Errorf("jwt secret key is not set")
	}
	return []byte(key), nil
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
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (service *UserService) createToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  user.Id,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Minute * 10).Unix(),
			"sub": strconv.FormatInt(user.Id, 10),
		})
	secretKey, err := service.getSecretKey()
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (service *UserService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	secretKey, err := service.getSecretKey()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
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

func (service *UserService) Login(ctx context.Context, name, password string) (token string, rawToken string, expiresAt time.Time, err error) { // returns jwt as string
	if err := service.validateCredentials(name, password); err != nil {
		return "", "", time.Time{}, err
	}
	user, err := service.Repo.GetUserByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", time.Time{}, ErrInvalidCredentials
		}
		return "", "", time.Time{}, err
	}
	if err := service.comparePasswords(password, user.PasswordHash); err != nil {
		return "", "", time.Time{}, ErrInvalidCredentials
	}
	token, err = service.createToken(&user)
	if err != nil {
		return "", "", time.Time{}, err
	}
	rawToken, expiresAt, err = service.IssueRefreshToken(ctx, user.Id)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return token, rawToken, expiresAt, nil
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

func (service *UserService) Logout(ctx context.Context, refreshTokenRaw string) error {
	if strings.TrimSpace(refreshTokenRaw) == "" {
		return ErrMissingRefreshToken
	}

	hash := service.HashRefreshToken(refreshTokenRaw)

	rt, err := service.Repo.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	if rt.RevokedAt != nil {
		return ErrInvalidRefreshToken
	}

	return service.Repo.Revoke(ctx, rt.Id, time.Now().UTC())
}

func (service *UserService) Refresh(ctx context.Context, refreshTokenRaw string) (newAccessToken string, newRefreshRaw string, newRefreshExpiresAt time.Time, err error) {
	if strings.TrimSpace(refreshTokenRaw) == "" {
		return "", "", time.Time{}, ErrMissingRefreshToken
	}

	now := time.Now().UTC()
	hash := service.HashRefreshToken(refreshTokenRaw)

	rt, err := service.Repo.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", time.Time{}, ErrInvalidRefreshToken
		}
		return "", "", time.Time{}, err
	}

	if rt.RevokedAt != nil {
		return "", "", time.Time{}, ErrInvalidRefreshToken
	}
	if !rt.ExpiresAt.After(now) {
		return "", "", time.Time{}, ErrInvalidRefreshToken
	}

	u := models.User{Id: rt.UserId}
	newAccessToken, err = service.createToken(&u)
	if err != nil {
		return "", "", time.Time{}, err
	}

	newRefreshRaw, err = service.GenerateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, err
	}
	newHash := service.HashRefreshToken(newRefreshRaw)
	newRefreshExpiresAt = now.Add(time.Hour * 24 * 7)

	_, err = service.Repo.Rotate(ctx, rt.Id, rt.UserId, newHash, newRefreshExpiresAt, now)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return newAccessToken, newRefreshRaw, newRefreshExpiresAt, nil
}

func (service *UserService) GenerateRefreshToken() (string, error) {
	const nBytes = 32

	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (service *UserService) HashRefreshToken(refreshToken string) string {
	sum := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(sum[:])
}

func (service *UserService) IssueRefreshToken(ctx context.Context, userID int64) (string, time.Time, error) {
	rawToken, err := service.GenerateRefreshToken()
	if err != nil {
		return "", time.Time{}, err
	}

	hash := service.HashRefreshToken(rawToken)
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 7)

	_, err = service.Repo.InsertRefreshToken(ctx, userID, hash, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return rawToken, expiresAt, nil
}
