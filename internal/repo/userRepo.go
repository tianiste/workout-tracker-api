package repo

import (
	"database/sql"
	"workout-tracker/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (repo *UserRepo) GetUserByName(name string) (models.User, error) {
	var user models.User
	if err := repo.DB.QueryRow("SELECT id, name, pass_hash FROM users WHERE name = ?", name).Scan(&user.Id, &user.Username, &user.PasswordHash); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (repo *UserRepo) InsertUser(name, passwordHash string) (int64, error) {
	res, err := repo.DB.Exec("INSERT INTO users (name, pass_hash), VALUES(?, ?)", name, passwordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
