package oauth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arinji2/vocab-thing/internal/database"
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
		},
		Db: db,
	}
}

// FetchAuthUser returns an AuthUser instance for a guest from DB.
func (p *Guest) FetchGuestUser() (*models.User, error) {
	var username string
	totalRuns := 0
	for {
		if totalRuns > 5 {
			return nil, fmt.Errorf("exceeding 5 total runs for generating guestID")
		}
		totalRuns++
		randomID := idgen.GenerateRandomID(6, idgen.NumberCharset)

		guestID := fmt.Sprintf("Guest-%s", randomID)
		var guestUser database.UserModel
		_, err := guestUser.ByUsername(p.Provider.Ctx, guestID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				username = guestID
				break
			} else {
				return nil, fmt.Errorf("error with checking guest unique username: %w", err)
			}
		}

	}

	user := &models.User{
		Username: username,
	}

	return user, nil
}
