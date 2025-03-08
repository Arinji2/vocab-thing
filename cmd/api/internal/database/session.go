package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arinji2/vocab-thing/internal/models"
)

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Create(ctx context.Context, session *models.Session) error {
	fmt.Println("Creating session")
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	session.CreatedAt = time.Now().UTC()
	query := `INSERT INTO sessions (id, userId, providerId, fingerprint, ip, expiresAt)
          VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, session.UserID, session.ProviderID, session.Fingerprint, session.IP, session.ExpiresAt.Format(time.RFC3339)).Scan(&session.ID)
	if err != nil {
		return fmt.Errorf("error with session creation of userID %s: %w", session.ID, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
