package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/errorcode"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils"
)

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Create(ctx context.Context, session *models.Session) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return errorcode.ErrTransactionStart
	}
	defer tx.Rollback()

	session.CreatedAt = time.Now().UTC()
	query := `INSERT INTO sessions (id, userId, providerId, fingerprint, ip, expiresAt)
          VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, session.UserID, session.ProviderID, session.Fingerprint, session.IP, session.ExpiresAt.Format(time.RFC3339)).Scan(&session.ID)
	if err != nil {
		log.Printf("error with session creation of userID %s: %s", session.ID, err.Error())
		return errorcode.ErrSessionCreate
	}
	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return errorcode.ErrTransactionCommit
	}

	return nil
}

func (m *SessionModel) ByUserID(ctx context.Context, id string) ([]models.Session, error) {
	query := `SELECT id, userId, providerId, fingerprint, ip, expiresAt, createdAt FROM sessions WHERE userId = ?`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("querying sessions: %s", err.Error())
		return nil, errorcode.ErrSessionQuery
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
			log.Printf("scanning session row: %s", err.Error())
			return nil, errorcode.ErrScanningRow
		}

		session.CreatedAt, _ = utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for session %s", createdAtStr, session.ID))
		session.ExpiresAt, _ = utils.StringToTime(expiresAtStr, fmt.Sprintf("Warning: could not parse expiresAt '%s' for session %s", expiresAtStr, session.ID))

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterating provider rows: %s", err.Error())
		return nil, errorcode.ErrIteratingRows
	}

	if len(sessions) == 0 {
		log.Printf("no sessions found for userID %s: %s", id, sql.ErrNoRows)
		return nil, sql.ErrNoRows
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
		log.Printf("querying sessions with provider: %s", err.Error())
		return nil, errorcode.ErrSessionQuery
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
			log.Printf("scanning session row: %s", err.Error())
			return nil, errorcode.ErrScanningRow
		}

		session.CreatedAt, _ = utils.StringToTime(createdAtStr, "")
		session.ExpiresAt, _ = utils.StringToTime(expiresAtStr, "")
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterating session rows: %s", err.Error())
		return nil, errorcode.ErrIteratingRows
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
			return models.Session{}, errorcode.ErrNoSession
		}
		log.Printf("scanning session row: %s", err.Error())
		return models.Session{}, errorcode.ErrScanningRow
	}

	parsedTime, err := utils.StringToTime(expiresAtStr, fmt.Sprintf("Warning: could not parse expiresAt '%s' for session %s", expiresAtStr, sessionID))
	if err != nil {
		log.Printf("could not parse expiresAt '%s' for session %s: %s", expiresAtStr, sessionID, err.Error())
	}
	session.ExpiresAt = parsedTime
	if session.ExpiresAt.Before(time.Now()) {
		return models.Session{}, auth.ErrSessionExpired
	}
	session.ExpiresAt = parsedTime

	return session, nil
}
