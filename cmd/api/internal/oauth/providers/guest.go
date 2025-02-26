package providers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arinji2/vocab-thing/internal/database/users"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils/idgen"
)

type Guest struct {
	Provider BaseProvider
	Db       *sql.DB
}

func NewGuestProvider(db *sql.DB) *Guest {
	return &Guest{
		Provider: BaseProvider{
			ProviderType: "guest",
			Ctx:          context.Background(),
			AuthURL:      "",
			TokenURL:     "",
			UserInfoURL:  "",
			Scopes:       []string{},
			ClientId:     "",
			ClientSecret: "",
			RedirectURL:  "",
		},
		Db: db,
	}
}

// FetchAuthUser returns an AuthUser instance for a guest from DB.
func (p *Guest) FetchGuestUser() (*models.AuthUser, error) {
	var username string
	totalRuns := 0
	for {
		if totalRuns > 5 {
			return nil, fmt.Errorf("exceeding 5 total runs for generating guestID")
		}
		totalRuns++
		randomID, err := idgen.GenerateRandomID(6, idgen.NumberCharset)
		if err != nil {
			return nil, fmt.Errorf("error with generating guest id: %w", err)
		}

		guestID := fmt.Sprintf("Guest-%s", randomID)
		var guestUser users.UserModel
		_, err = guestUser.ByUsername(p.Provider.Ctx, guestID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				username = guestID
				break
			} else {
				return nil, fmt.Errorf("error with checking guest unique username: %w", err)
			}
		}

	}

	id, err := idgen.GenerateRandomID(idgen.DefaultIDSize, idgen.URLSafeAlphanumericCharset)
	if err != nil {
		return nil, err
	}

	user := &models.AuthUser{
		Type:     p.Provider.ProviderType,
		Id:       id,
		Username: username,
	}

	return user, nil
}
