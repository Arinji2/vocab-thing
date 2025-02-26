package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/tools/types"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) All(ctx context.Context) ([]models.User, error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT id, username, email, createdAt FROM users;`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var createdAtStr string

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("scanning user row: %w", err)
		}

		parsedTime, err := types.ParseDateTime(createdAtStr)
		if err != nil {
			fmt.Printf("Warning: could not parse createdAt '%s' for user %d: %v\n", createdAtStr, user.ID, err)
			user.CreatedAt = types.DateTime{}
		} else {
			user.CreatedAt = parsedTime
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return users, nil
}

func (m *UserModel) ByID(ctx context.Context, id string) (models.User, error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT id, username, email, createdAt FROM users WHERE id = ?`

	row := tx.QueryRowContext(ctx, query, id)

	var user models.User
	var createdAtStr string

	err = row.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found with id %s: %w", id, err)
		}
		return models.User{}, fmt.Errorf("scanning user row: %w", err)
	}

	parsedTime, err := types.ParseDateTime(createdAtStr)
	if err != nil {
		fmt.Printf("Warning: could not parse createdAt '%s' for user %d: %v\n", createdAtStr, user.ID, err)
		user.CreatedAt = types.DateTime{}
	} else {
		user.CreatedAt = parsedTime
	}

	if err := tx.Commit(); err != nil {
		return models.User{}, fmt.Errorf("committing transaction: %w", err)
	}

	return user, nil
}

func (m *UserModel) ByUsername(ctx context.Context, username string) (models.User, error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT id, username, email, createdAt FROM users WHERE username = ?`

	row := tx.QueryRowContext(ctx, query, username)

	var user models.User
	var createdAtStr string

	err = row.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found with username %s: %w", username, err)
		}
		return models.User{}, fmt.Errorf("scanning user row: %w", err)
	}

	parsedTime, err := types.ParseDateTime(createdAtStr)
	if err != nil {
		fmt.Printf("Warning: could not parse createdAt '%s' for user %d: %v\n", createdAtStr, user.ID, err)
		user.CreatedAt = types.DateTime{}
	} else {
		user.CreatedAt = parsedTime
	}

	if err := tx.Commit(); err != nil {
		return models.User{}, fmt.Errorf("committing transaction: %w", err)
	}

	return user, nil
}
