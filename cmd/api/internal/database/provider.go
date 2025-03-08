package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils"
)

type ProviderModel struct {
	DB *sql.DB
}

func (p *ProviderModel) Create(ctx context.Context, provider *models.OauthProvider) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()
	query := `INSERT INTO providers (id, userId, type, refreshToken, accessToken, expiresAt)
          VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, provider.UserID, provider.Type, provider.RefreshToken, provider.AccessToken, provider.ExpiresAt.Format(time.RFC3339)).Scan(&provider.ID)
	if err != nil {
		return fmt.Errorf("error with provider creation of userID %s and provider type %s: %w", provider.UserID, provider.Type, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *ProviderModel) ByUserID(ctx context.Context, id string) ([]models.OauthProvider, error) {
	query := `SELECT id, userID, type, accessToken, expiresAt, refreshToken, createdAt FROM providers WHERE userID = ?`

	rows, err := p.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("querying providers: %w", err)
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
			return nil, fmt.Errorf("scanning provider row: %w", err)
		}

		provider.CreatedAt, _ = utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for provider %s", createdAtStr, provider.ID))
		provider.ExpiresAt, _ = utils.StringToTime(expiresAtStr, fmt.Sprintf("Warning: could not parse expiresAt '%s' for provider %s", expiresAtStr, provider.ID))

		providers = append(providers, provider)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating provider rows: %w", err)
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers found for userID %s: %w", id, sql.ErrNoRows)
	}

	return providers, nil
}
