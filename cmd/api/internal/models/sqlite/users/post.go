package users

import (
	"context"
	"fmt"

	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/tools/types"
)

func (m *UserModel) Create(ctx context.Context, user models.User) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	user.CreatedAt = types.NowDateTime()
	query := `INSERT INTO users (id, username, email, createdAt) VALUES (?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.CreatedAt.String())
	if err != nil {
		return fmt.Errorf("error with user creation of username %s, userID %d: %w", user.Username, user.ID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
