package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/arinji2/vocab-thing/internal/errorcode"
	"github.com/arinji2/vocab-thing/internal/models"
)

type SyncModel struct {
	DB *sql.DB
}

func (s *SyncModel) CreateSync(ctx context.Context, userID string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return errorcode.ErrTransactionStart
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
		log.Printf("error creating sync record for userID %s: %s", syncData.UserID, err.Error())
		return errorcode.ErrDBCreate
	}

	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return errorcode.ErrTransactionCommit
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
		log.Printf("error fetching sync record for userID %s: %s", userID, err.Error())
		return nil, errorcode.ErrDBQuery
	}
	return &syncData, nil
}

func (s *SyncModel) ManualSync(ctx context.Context, userID string) (*models.SyncMetadata, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return nil, errorcode.ErrTransactionStart
	}
	defer tx.Rollback()
	existingSync, err := s.ByUserID(ctx, userID)
	if err != nil {
		log.Printf("error fetching sync record for userID %s: %s", userID, err.Error())
		return nil, errorcode.ErrDBQuery
	}

	if existingSync.LastUpdatedAt.After(time.Now().Add(-30 * time.Minute)) {
		return existingSync, errorcode.ErrManualSyncLimit.WithDetails(map[string]int{"waitFor": int(time.Until(existingSync.LastUpdatedAt.Add(30 * time.Minute)).Minutes())})
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
		log.Printf("error updating sync record for userID %s: %s", syncData.UserID, err.Error())
		return nil, errorcode.ErrDBUpdate
	}

	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return nil, errorcode.ErrTransactionCommit
	}
	return &syncData, nil
}
