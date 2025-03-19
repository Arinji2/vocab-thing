package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/arinji2/vocab-thing/internal/errorcode"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils"
)

type ProviderModel struct {
	DB *sql.DB
}

func (p *ProviderModel) Create(ctx context.Context, provider *models.OauthProvider) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return errorcode.ErrTransactionStart
	}
	defer tx.Rollback()
	query := `INSERT INTO providers (id, userId, type, refreshToken, accessToken, expiresAt)
          VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, provider.UserID, provider.Type, provider.RefreshToken, provider.AccessToken, provider.ExpiresAt.Format(time.RFC3339)).Scan(&provider.ID)
	if err != nil {
		log.Printf("error with provider creation of userID %s and provider type %s: %s", provider.UserID, provider.Type, err.Error())
		return errorcode.ErrDBCreate

	}
	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return errorcode.ErrTransactionCommit
	}

	return nil
}

func (p *ProviderModel) ByUserID(ctx context.Context, id string) ([]models.OauthProvider, error) {
	query := `SELECT id, userID, type, accessToken, expiresAt, refreshToken, createdAt FROM providers WHERE userID = ?`

	rows, err := p.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("querying providers: %s", err.Error())
		return nil, errorcode.ErrDBQuery
	}
	defer rows.Close()

	var providers []models.OauthProvider

	for rows.Next() {
		var provider models.OauthProvider
		var createdAtStr, expiresAtStr string

		err := rows.Scan(
			&provider.ID,
			&provider.UserID,
			&provider.Type,
			&provider.AccessToken,
			&expiresAtStr,
			&provider.RefreshToken,
			&createdAtStr,
		)
		if err != nil {
			log.Printf("scanning provider row: %s", err.Error())
			return nil, errorcode.ErrScanningRow
		}

		provider.CreatedAt, err = utils.StringToTime(createdAtStr)
		if err != nil {
			log.Printf("Warning: could not parse createdAt '%s' for provider %s", createdAtStr, provider.ID)
			provider.CreatedAt = time.Now().UTC()
		}

		provider.ExpiresAt, err = utils.StringToTime(expiresAtStr)
		if err != nil {
			log.Printf("Warning: could not parse expiresAt '%s' for provider %s", expiresAtStr, provider.ID)
			provider.ExpiresAt = time.Now().UTC()
		}

		providers = append(providers, provider)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterating provider rows: %s", err.Error())
		return nil, errorcode.ErrIteratingRows
	}

	if len(providers) == 0 {
		log.Printf("no providers found for userID %s: %s", id, sql.ErrNoRows)
		return nil, sql.ErrNoRows
	}

	return providers, nil
}
