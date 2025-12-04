package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	Id        int       `json:"id"`
	OwnerId   string    `json:"owner_id"`
	Hash      string    `json:"hash"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshTokenQueries provides database operations for refresh tokens
type RefreshTokenQueries struct {
	db *DB
}

// NewRefreshTokenQueries creates a new RefreshTokenQueries instance
func NewRefreshTokenQueries(db *DB) *RefreshTokenQueries {
	return &RefreshTokenQueries{db: db}
}

// Create inserts a new user
func (q *RefreshTokenQueries) Create(ownerId, tokenHash string) (*RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (owner_id, hash, expires_at)
		VALUES (?, ?, datetime('now', '+30 days'))
	`

	if ownerId == "" {
		return nil, fmt.Errorf("owner_id cannot be empty")
	}

	if tokenHash == "" {
		return nil, fmt.Errorf("token_hash cannot be empty")
	}

	result, err := q.db.Exec(query, ownerId, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token for user '%s': %s", ownerId, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return q.GetByID(int(id))
}

// GetByID retrieves valid refresh tokens for a specific user
func (q *RefreshTokenQueries) GetByID(tokenId int) (*RefreshToken, error) {
	query := `
		SELECT id, owner_id, hash, issued_at, expires_at
		FROM refresh_tokens
		WHERE owner_id = ?
	`

	var token RefreshToken
	err := q.db.QueryRow(query, tokenId).Scan(
		&token.Id,
		&token.OwnerId,
		&token.Hash,
		&token.IssuedAt,
		&token.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token with id '%d' not found", tokenId)
		}

		return nil, fmt.Errorf("failed to get refresh token with id '%d': %w", tokenId, err)
	}

	return &token, nil
}

// GetValidByUserID retrieves valid refresh tokens for a specific user
func (q *RefreshTokenQueries) GetValidByUserID(userId int) ([]RefreshToken, error) {
	query := `
		SELECT id, owner_id, hash, issued_at, expires_at
		FROM refresh_tokens
		WHERE owner_id = ? AND expires_at > datetime('now')
	`

	rows, err := q.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid refresh tokens for user '%d': %w", userId, err)
	}

	var tokens []RefreshToken
	defer rows.Close()

	for rows.Next() {
		token := RefreshToken{}
		err = rows.Scan(&token.Id, &token.OwnerId, &token.Hash, &token.IssuedAt, &token.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan refresh token: %w", err)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// Count returns the total number of refresh tokens
func (q *RefreshTokenQueries) Count() (int, error) {
	query := "SELECT COUNT(*) FROM refresh_tokens"

	var count int
	err := q.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count refresh tokens: %w", err)
	}

	return count, nil
}
