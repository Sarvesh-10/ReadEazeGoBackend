package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type RefreshTokenRepository interface {
	Store(userID int, tokenHash string, expiresAt time.Time) error
	Validate(tokenHash string) (userID int, emailID string, err error)
	Revoke(tokenHash string) error
	RevokeAll(userID int) error
}

type RefreshTokenRepositoryImpl struct {
	db     *sql.DB
	logger *utility.Logger
}

func NewRefreshTokenRepository(db *sql.DB, logger *utility.Logger) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{db, logger}
}
func (r *RefreshTokenRepositoryImpl) Store(userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)", userID, tokenHash, expiresAt)
	if err != nil {
		r.logger.Error("Error storing refresh token: %s", err.Error())
		return err
	}
	return nil
}
func (r *RefreshTokenRepositoryImpl) Validate(tokenHash string) (int, string, error) {
	var userID int
	var emailID string
	var revoked bool
	var expiresAt time.Time

	err := r.db.QueryRow(`
	SELECT u.id, u.email, rt.revoked, rt.expires_at
	FROM refresh_tokens rt
	JOIN users u ON u.id = rt.user_id
	WHERE rt.token_hash = $1
`, tokenHash).Scan(&userID, &emailID, &revoked, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New("refresh token not found")
		}
		r.logger.Error("Error validating refresh token: %s", err.Error())
		return 0, "", err
	}

	if revoked {
		return 0, "", errors.New("token has been revoked")
	}
	if time.Now().After(expiresAt) {
		return 0, "", errors.New("token has expired")
	}

	return userID, emailID, nil
}

func (r *RefreshTokenRepositoryImpl) Revoke(tokenHash string) error {
	_, err := r.db.Exec(`
		UPDATE refresh_tokens SET revoked = TRUE
		WHERE token_hash = $1
	`, tokenHash)

	if err != nil {
		r.logger.Error("Error revoking refresh token: %s", err.Error())
	}
	return err
}

func (r *RefreshTokenRepositoryImpl) RevokeAll(userID int) error {
	_, err := r.db.Exec(`
		UPDATE refresh_tokens SET revoked = TRUE
		WHERE user_id = $1
	`, userID)

	if err != nil {
		r.logger.Error("Error revoking all tokens for user %d: %s", userID, err.Error())
	}
	return err
}
