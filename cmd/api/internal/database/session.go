package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils"
)

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Create(ctx context.Context, session *models.Session) error {
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

func (m *SessionModel) ByUserID(ctx context.Context, id string) ([]models.Session, error) {
	query := `SELECT id, userId, providerId, fingerprint, ip, expiresAt, createdAt FROM sessions WHERE userId = ?`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("querying sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session

	for rows.Next() {
		var session models.Session
		var createdAtStr, expiresAtStr string

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.ProviderID,
			&session.Fingerprint,
			&session.IP,
			&expiresAtStr,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning session row: %w", err)
		}

		session.CreatedAt, _ = utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for session %s", createdAtStr, session.ID))
		session.ExpiresAt, _ = utils.StringToTime(expiresAtStr, fmt.Sprintf("Warning: could not parse expiresAt '%s' for session %s", expiresAtStr, session.ID))

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating provider rows: %w", err)
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found for userID %s: %w", id, sql.ErrNoRows)
	}

	return sessions, nil
}

func (m *SessionModel) ByUserIDWithProvider(ctx context.Context, id string) ([]models.Session, error) {
	query := `SELECT s.id, s.userId, s.providerId, p.type, s.fingerprint, s.ip, s.expiresAt, s.createdAt
	          FROM sessions s
	          JOIN providers p ON s.providerId = p.id
	          WHERE s.userId = ?`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("querying sessions with provider: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		var createdAtStr, expiresAtStr string

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.ProviderID,
			&session.ProviderType,
			&session.Fingerprint,
			&session.IP,
			&expiresAtStr,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning session row: %w", err)
		}

		session.CreatedAt, _ = utils.StringToTime(createdAtStr, "")
		session.ExpiresAt, _ = utils.StringToTime(expiresAtStr, "")
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating session rows: %w", err)
	}

	if len(sessions) == 0 {
		return nil, sql.ErrNoRows
	}

	return sessions, nil
}

func (m *SessionModel) Validate(ctx context.Context, sessionID string) (models.Session, error) {
	query := `SELECT id, userId, expiresAt FROM sessions WHERE id = ?`

	row := m.DB.QueryRowContext(ctx, query, sessionID)

	var session models.Session
	var expiresAtStr string

	err := row.Scan(&session.ID, &session.UserID, &expiresAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Session{}, fmt.Errorf("session not found with id %s: %w", sessionID, err)
		}
		return models.Session{}, fmt.Errorf("scanning session row: %w", err)
	}

	parsedTime, err := utils.StringToTime(expiresAtStr, fmt.Sprintf("Warning: could not parse expiresAt '%s' for session %s", expiresAtStr, sessionID))
	if err != nil {
		fmt.Println(err.Error())
	}
	session.ExpiresAt = parsedTime
	if session.ExpiresAt.Before(time.Now()) {
		return models.Session{}, auth.ErrSessionExpired
	}
	session.ExpiresAt = parsedTime

	return session, nil
}
