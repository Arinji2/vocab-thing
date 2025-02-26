package providers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arinji2/vocab-thing/internal/models/idgen"
	"github.com/arinji2/vocab-thing/internal/models/sqlite"
)

type Guest struct {
	BaseProvider
	Db *sql.DB
}

func NewGuestProvider(db *sql.DB) *Google {
	return &Google{
		BaseProvider{
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
	}
}

// FetchAuthUser returns an AuthUser instance for a guest from DB.
func (p *Guest) FetchGuestUser() (*AuthUser, error) {
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
		var guestUser sqlite.UserModel
		_, err = guestUser.ByUsername(p.Ctx, guestID)
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

	user := &AuthUser{
		Type:     p.ProviderType,
		Id:       id,
		Username: username,
	}

	return user, nil
}
