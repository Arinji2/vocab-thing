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

func (s *SyncModel) ByUserID(ctx context.Context, userID string) (*models.SyncMetadata, error) {
	query := `
		SELECT id, userId, lastUpdatedAt
		FROM sync_metadata
		WHERE userId = ?
	`
	var syncData models.SyncMetadata
	err := s.DB.QueryRowContext(ctx, query, userID).Scan(&syncData.ID, &syncData.UserID, &syncData.LastUpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching sync record for userID %s: %w", userID, err)
	}
	return &syncData, nil
}

func (s *SyncModel) ManualSync(ctx context.Context, userID string) (*models.SyncMetadata, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()
	existingSync, err := s.ByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching sync record for userID %s: %w", userID, err)
	}

	if existingSync.LastUpdatedAt.After(time.Now().Add(-30 * time.Minute)) {
		return existingSync, fmt.Errorf("sync record for userID %s is already up to date, please wait %v minutes before resyncing", userID, int(time.Until(existingSync.LastUpdatedAt.Add(30*time.Minute)).Minutes()))
	}

	syncData := models.SyncMetadata{
		UserID:        userID,
		LastUpdatedAt: time.Now().UTC(),
	}

	query := `
    UPDATE sync_metadata SET lastUpdatedAt = ? WHERE userId = ?
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query, syncData.LastUpdatedAt, syncData.UserID).Scan(&syncData.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating sync record for userID %s: %w", userID, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}
	return &syncData, nil
}
