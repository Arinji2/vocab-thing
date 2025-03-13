package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils"
)

type PhraseModel struct {
	DB *sql.DB
}

func (p *PhraseModel) CreatePhrase(ctx context.Context, phrase *models.Phrase) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	phrase.CreatedAt = time.Now().UTC()
	query := `
            INSERT INTO phrases (id, userId, phrase, phraseDefinition, pinned, foundIn, public, usageCount, createdAt)
            VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?, ?, ?, ?, ?) 
  RETURNING id
            `

	err = tx.QueryRowContext(ctx, query, phrase.UserID, phrase.Phrase, phrase.PhraseDefinition, phrase.Pinned, phrase.FoundIn, phrase.Public, phrase.UsageCount, phrase.CreatedAt.Format(time.RFC3339)).Scan(&phrase.ID)
	if err != nil {
		return fmt.Errorf("error with phease creation of userID %s: %w", phrase.UserID, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *PhraseModel) CreateTag(ctx context.Context, phrase *models.PhraseTag) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	phrase.CreatedAt = time.Now().UTC()
	query := `
            INSERT INTO phrase_tags (id, phraseId, tagName, tagColor, createdAt)
            VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?) 
            RETURNING id
            `

	err = tx.QueryRowContext(ctx, query, phrase.PhraseID, phrase.TagName, phrase.TagColor, phrase.CreatedAt.Format(time.RFC3339)).Scan(&phrase.ID)
	if err != nil {
		return fmt.Errorf("error with phease tag creation of phraseID %s: %w", phrase.PhraseID, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *PhraseModel) ByID(ctx context.Context, id string, userID string) (*models.TaggedPhrase, error) {
	query := `
        SELECT p.id, p.userId, p.phrase, p.phraseDefinition, p.pinned, p.foundIn, p.public, p.usageCount, p.createdAt, 
        pt.id, pt.phraseId, pt.tagName, pt.tagColor, pt.createdAt
        FROM phrases p
        LEFT JOIN phraseTags pt ON p.id = pt.phraseId
        WHERE p.id = ? AND p.userId = ?
    `
	rows, err := p.DB.QueryContext(ctx, query, id, userID)
	if err != nil {
		return nil, fmt.Errorf("querying phrase: %w", err)
	}
	defer rows.Close()

	var taggedPhrase models.TaggedPhrase
	phraseLoaded := false

	for rows.Next() {
		var phraseCreatedAtStr, tagCreatedAtStr string
		var tag models.PhraseTag
		var tagID, phraseID sql.NullString

		err = rows.Scan(
			&taggedPhrase.Phrase.ID,
			&taggedPhrase.Phrase.UserID,
			&taggedPhrase.Phrase.Phrase,
			&taggedPhrase.Phrase.PhraseDefinition,
			&taggedPhrase.Phrase.Pinned,
			&taggedPhrase.Phrase.FoundIn,
			&taggedPhrase.Phrase.Public,
			&taggedPhrase.Phrase.UsageCount,
			&phraseCreatedAtStr,
			&tagID,
			&phraseID,
			&tag.TagName,
			&tag.TagColor,
			&tagCreatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning taggedPhrase row: %w", err)
		}

		// Parse the phrase creation time on the first iteration
		if !phraseLoaded {
			taggedPhrase.Phrase.CreatedAt, _ = utils.StringToTime(
				phraseCreatedAtStr,
				fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase %s",
					phraseCreatedAtStr, taggedPhrase.Phrase.Phrase),
			)
			phraseLoaded = true
		}

		// Only add the tag if it's not NULL (some phrases might not have tags)
		if tagID.Valid {
			tag.ID = tagID.String
			tag.PhraseID = phraseID.String
			tag.CreatedAt, _ = utils.StringToTime(
				tagCreatedAtStr,
				fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase tag %s for phrase %s",
					tagCreatedAtStr, tag.TagName, taggedPhrase.Phrase.Phrase),
			)

			taggedPhrase.Tag = append(taggedPhrase.Tag, tag)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating phrase rows: %w", err)
	}

	if !phraseLoaded {
		return nil, fmt.Errorf("phrase with id %s not found", id)
	}

	return &taggedPhrase, nil
}

func (p *PhraseModel) Search(ctx context.Context, searchQuery string, userID string) ([]models.TaggedPhrase, error) {
	query := `
    WITH search_results AS (
        SELECT rowid FROM phrases_fts 
        WHERE phrase MATCH ? 
        ORDER BY rank
    )
    SELECT p.id, p.userId, p.phrase, p.phraseDefinition, p.pinned, p.foundIn, p.public, p.usageCount, p.createdAt, 
           pt.id, pt.phraseId, pt.tagName, pt.tagColor, pt.createdAt
    FROM search_results sr
    JOIN phrases p ON sr.rowid = p.id
    LEFT JOIN phrase_tags pt ON p.id = pt.phraseId
    WHERE p.userId = ?  
    ORDER BY p.usageCount DESC, p.createdAt DESC
`

	rows, err := p.DB.QueryContext(ctx, query, searchQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("searching phrases: %w", err)
	}
	defer rows.Close()

	uniquePhraseMap := make(map[string]*models.TaggedPhrase)

	for rows.Next() {
		var phraseCreatedAtStr, tagCreatedAtStr string
		var tagID, tagPhraseID, tagName, tagColor sql.NullString
		var phrase models.Phrase

		err = rows.Scan(
			&phrase.ID,
			&phrase.UserID,
			&phrase.Phrase,
			&phrase.PhraseDefinition,
			&phrase.Pinned,
			&phrase.FoundIn,
			&phrase.Public,
			&phrase.UsageCount,
			&phraseCreatedAtStr,
			&tagID,
			&tagPhraseID,
			&tagName,
			&tagColor,
			&tagCreatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning search result row: %w", err)
		}

		phrase.CreatedAt, _ = utils.StringToTime(
			phraseCreatedAtStr,
			fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase %s",
				phraseCreatedAtStr, phrase.Phrase),
		)

		if _, exists := uniquePhraseMap[phrase.ID]; !exists {
			uniquePhraseMap[phrase.ID] = &models.TaggedPhrase{
				Phrase: phrase,
				Tag:    []models.PhraseTag{},
			}
		}

		if tagID.Valid && tagName.Valid {
			var tag models.PhraseTag
			tag.ID = tagID.String
			tag.PhraseID = tagPhraseID.String
			tag.TagName = tagName.String
			tag.TagColor = tagColor.String

			tag.CreatedAt, _ = utils.StringToTime(
				tagCreatedAtStr,
				fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase tag %s for phrase %s",
					tagCreatedAtStr, tag.TagName, phrase.Phrase),
			)

			uniquePhraseMap[phrase.ID].Tag = append(uniquePhraseMap[phrase.ID].Tag, tag)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating search result rows: %w", err)
	}

	results := make([]models.TaggedPhrase, 0, len(uniquePhraseMap))
	for _, taggedPhrase := range uniquePhraseMap {
		results = append(results, *taggedPhrase)
	}

	return results, nil
}
