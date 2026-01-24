package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
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
	res, err := repo.DB.Exec("INSERT INTO users (name, pass_hash) VALUES(?, ?)", name, passwordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepo) InsertRefreshToken(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) (models.RefreshToken, error) {

	if userID <= 0 {
		return models.RefreshToken{}, errors.New("userID must be positive")
	}
	if tokenHash == "" {
		return models.RefreshToken{}, errors.New("tokenHash must not be empty")
	}

	res, err := r.DB.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, userID, tokenHash, expiresAt.UTC())
	if err != nil {
		return models.RefreshToken{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.RefreshToken{}, err
	}

	return r.getByID(ctx, id)
}

func (r *UserRepo) GetByHash(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	if tokenHash == "" {
		return models.RefreshToken{}, errors.New("refreshTokenHash must not be empty")
	}

	row := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, created_at, last_used_at, expires_at, revoked_at, replaced_by_token_id
		FROM refresh_tokens
		WHERE token_hash = ?
		LIMIT 1
	`, tokenHash)

	var rt models.RefreshToken
	var lastUsed sql.NullTime
	var revoked sql.NullTime
	var replacedBy sql.NullInt64

	err := row.Scan(
		&rt.Id,
		&rt.UserId,
		&rt.Hash,
		&rt.CreatedAt,
		&lastUsed,
		&rt.ExpiresAt,
		&revoked,
		&replacedBy,
	)
	if err != nil {
		return models.RefreshToken{}, err
	}

	if lastUsed.Valid {
		t := lastUsed.Time
		rt.LastUsedAt = &t
	}
	if revoked.Valid {
		t := revoked.Time
		rt.RevokedAt = &t
	}
	if replacedBy.Valid {
		v := replacedBy.Int64
		rt.ReplacedByTokenId = &v
	}

	return rt, nil
}

func (r *UserRepo) Revoke(ctx context.Context, tokenID int64, revokedAt time.Time) error {
	if tokenID <= 0 {
		return errors.New("tokenID must be positive")
	}

	_, err := r.DB.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = ?
		WHERE id = ? AND revoked_at IS NULL
	`, revokedAt.UTC(), tokenID)
	return err
}

func (r *UserRepo) Rotate(
	ctx context.Context,
	oldTokenID int64,
	newUserID int64,
	newHash string,
	newExpiresAt time.Time,
	rotatedAt time.Time,
) (models.RefreshToken, error) {

	if oldTokenID <= 0 {
		return models.RefreshToken{}, errors.New("oldTokenID must be positive")
	}
	if newUserID <= 0 {
		return models.RefreshToken{}, errors.New("newUserID must be positive")
	}
	if newHash == "" {
		return models.RefreshToken{}, errors.New("newHash must not be empty")
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.RefreshToken{}, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, newUserID, newHash, newExpiresAt.UTC())
	if err != nil {
		return models.RefreshToken{}, err
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return models.RefreshToken{}, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = ?, replaced_by_token_id = ?
		WHERE id = ? AND revoked_at IS NULL
	`, rotatedAt.UTC(), newID, oldTokenID)
	if err != nil {
		return models.RefreshToken{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.RefreshToken{}, err
	}

	return r.getByID(ctx, newID)
}

func (r *UserRepo) getByID(ctx context.Context, id int64) (models.RefreshToken, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, created_at, last_used_at, expires_at, revoked_at, replaced_by_token_id
		FROM refresh_tokens
		WHERE id = ?
		LIMIT 1
	`, id)

	var rt models.RefreshToken
	var lastUsed sql.NullTime
	var revoked sql.NullTime
	var replacedBy sql.NullInt64

	err := row.Scan(
		&rt.Id,
		&rt.UserId,
		&rt.Hash,
		&rt.CreatedAt,
		&lastUsed,
		&rt.ExpiresAt,
		&revoked,
		&replacedBy,
	)
	if err != nil {
		return models.RefreshToken{}, err
	}

	if lastUsed.Valid {
		t := lastUsed.Time
		rt.LastUsedAt = &t
	}
	if revoked.Valid {
		t := revoked.Time
		rt.RevokedAt = &t
	}
	if replacedBy.Valid {
		v := replacedBy.Int64
		rt.ReplacedByTokenId = &v
	}

	return rt, nil
}
