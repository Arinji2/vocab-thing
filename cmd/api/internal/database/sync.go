package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arinji2/vocab-thing/internal/models"
)

type SyncModel struct {
	DB *sql.DB
}

func (s *SyncModel) CreateSync(ctx context.Context, userID string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	syncData := models.SyncMetadata{
		UserID:        userID,
		LastUpdatedAt: time.Now().UTC(),
	}

	query := `
		INSERT INTO sync_metadata (id, userId, lastUpdatedAt)
		VALUES (lower(hex(randomblob(16))), ?, ?)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query, syncData.UserID, syncData.LastUpdatedAt).Scan(&syncData.ID)
	if err != nil {
		return fmt.Errorf("error creating sync record for userID %s: %w", userID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
