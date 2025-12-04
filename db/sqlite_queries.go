package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// User represents a user in the database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserQueries provides database operations for users
type UserQueries struct {
	db *DB
}

// NewUserQueries creates a new UserQueries instance
func NewUserQueries(db *DB) *UserQueries {
	return &UserQueries{db: db}
}

// Create inserts a new user
func (q *UserQueries) Create(email, name string) (*User, error) {
	query := `
		INSERT INTO users (email, name, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := q.db.Exec(query, email, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return q.GetByID(int(id))
}

// GetByID retrieves a user by ID
func (q *UserQueries) GetByID(id int) (*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := q.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (q *UserQueries) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user User
	err := q.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// List retrieves all users with optional limit and offset
func (q *UserQueries) List(limit, offset int) ([]User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := q.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

// Update modifies an existing user
func (q *UserQueries) Update(id int, email, name string) (*User, error) {
	query := `
		UPDATE users
		SET email = ?, name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := q.db.Exec(query, email, name, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return q.GetByID(id)
}

// Delete removes a user by ID
func (q *UserQueries) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := q.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Count returns the total number of users
func (q *UserQueries) Count() (int, error) {
	query := "SELECT COUNT(*) FROM users"

	var count int
	err := q.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
