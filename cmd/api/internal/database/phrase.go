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
		phrase, tag, err := scanTaggedPhrase(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning taggedPhrase row: %w", err)
		}

		if !phraseLoaded {
			taggedPhrase.Phrase = phrase
			phraseLoaded = true
		}

		if tag != nil {
			taggedPhrase.Tag = append(taggedPhrase.Tag, *tag)
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

type Scanner interface {
	Scan(dest ...any) error
}

func scanTaggedPhrase(scanner Scanner) (models.Phrase, *models.PhraseTag, error) {
	var phrase models.Phrase
	var phraseCreatedAtStr string
	var tagID, tagPhraseID, tagName, tagColor sql.NullString
	var tagCreatedAtStr string

	err := scanner.Scan(
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
		return phrase, nil, err
	}

	phrase.CreatedAt, _ = utils.StringToTime(
		phraseCreatedAtStr,
		fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase %s", phraseCreatedAtStr, phrase.Phrase),
	)

	if tagID.Valid && tagName.Valid {
		tag := models.PhraseTag{
			ID:       tagID.String,
			PhraseID: tagPhraseID.String,
			TagName:  tagName.String,
			TagColor: tagColor.String,
		}
		tag.CreatedAt, _ = utils.StringToTime(
			tagCreatedAtStr,
			fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase tag %s for phrase %s", tagCreatedAtStr, tag.TagName, phrase.Phrase),
		)
		return phrase, &tag, nil
	}
	return phrase, nil, nil
}
