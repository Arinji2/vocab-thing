package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arinji2/vocab-thing/internal/models"
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

		parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			fmt.Printf("Warning: could not parse createdAt '%s' for user %d: %v\n", createdAtStr, user.ID, err)
			user.CreatedAt = time.Time{}
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
