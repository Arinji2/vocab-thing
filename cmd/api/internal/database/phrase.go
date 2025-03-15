package database

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
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
        LEFT JOIN phrase_tags pt ON p.id = pt.phraseId
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

func (p *PhraseModel) CountTotalPages(ctx context.Context, pageSize int, userID string) (int, error) {
	var totalRecords int
	query := `SELECT COUNT(*) FROM phrases WHERE userId = ?`
	err := p.DB.QueryRowContext(ctx, query, userID).Scan(&totalRecords)
	if err != nil {
		return 0, fmt.Errorf("counting phrases: %w", err)
	}
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return totalPages, nil
}

func (p *PhraseModel) All(ctx context.Context, pageNumber, pageSize int, sortBy, order, groupBy string, userID string) ([]models.TaggedPhrase, error) {
	totalPages, err := p.CountTotalPages(ctx, pageSize, userID)
	if err != nil {
		return nil, fmt.Errorf("counting total pages: %w", err)
	}
	if pageNumber > totalPages {
		return nil, fmt.Errorf("page number %d is greater than total pages %d", pageNumber, totalPages)
	}

	phrasesQuery := fmt.Sprintf(`
        SELECT id, userId, phrase, phraseDefinition, pinned, foundIn, public, usageCount, createdAt
        FROM phrases
        WHERE userId = ?
        ORDER BY %s %s
        LIMIT ? OFFSET ?
    `, sortBy, order)
	offset := (pageNumber - 1) * pageSize
	phraseRows, err := p.DB.QueryContext(ctx, phrasesQuery, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("querying paginated phrases: %w", err)
	}
	defer phraseRows.Close()

	var phrases []models.TaggedPhrase
	phraseIDs := []string{}
	phraseIndex := make(map[string]int) // Map to track index of phrases in slice

	for phraseRows.Next() {
		var phrase models.Phrase
		var phraseCreatedAtStr string

		err = phraseRows.Scan(
			&phrase.ID,
			&phrase.UserID,
			&phrase.Phrase,
			&phrase.PhraseDefinition,
			&phrase.Pinned,
			&phrase.FoundIn,
			&phrase.Public,
			&phrase.UsageCount,
			&phraseCreatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning phrase row: %w", err)
		}

		phrase.CreatedAt, _ = utils.StringToTime(
			phraseCreatedAtStr,
			fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase %s", phraseCreatedAtStr, phrase.Phrase),
		)

		phraseIDs = append(phraseIDs, phrase.ID)
		phraseIndex[phrase.ID] = len(phrases)
		phrases = append(phrases, models.TaggedPhrase{Phrase: phrase, Tag: []models.PhraseTag{}})
	}

	if err := phraseRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating phrase rows: %w", err)
	}

	if len(phraseIDs) == 0 {
		return []models.TaggedPhrase{}, nil
	}

	placeholders := make([]string, len(phraseIDs))
	args := make([]any, len(phraseIDs))
	for i, id := range phraseIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	tagsQuery := fmt.Sprintf(`
        SELECT id, phraseId, tagName, tagColor, createdAt
        FROM phrase_tags
        WHERE phraseId IN (%s)
    `, strings.Join(placeholders, ","))

	tagRows, err := p.DB.QueryContext(ctx, tagsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("querying tags: %w", err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tag models.PhraseTag
		var tagCreatedAtStr string

		err = tagRows.Scan(&tag.ID, &tag.PhraseID, &tag.TagName, &tag.TagColor, &tagCreatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("scanning tag row: %w", err)
		}

		tag.CreatedAt, _ = utils.StringToTime(
			tagCreatedAtStr,
			fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase tag %s", tagCreatedAtStr, tag.TagName),
		)

		if idx, exists := phraseIndex[tag.PhraseID]; exists {
			phrases[idx].Tag = append(phrases[idx].Tag, tag)
		}
	}

	if err := tagRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tag rows: %w", err)
	}
	if groupBy == "foundIn" {
		grouped := make(map[string][]models.TaggedPhrase)
		orderedGroups := []string{}

		for _, phrase := range phrases {
			if _, exists := grouped[phrase.Phrase.FoundIn]; !exists {
				orderedGroups = append(orderedGroups, phrase.Phrase.FoundIn)
			}
			grouped[phrase.Phrase.FoundIn] = append(grouped[phrase.Phrase.FoundIn], phrase)
		}

		var groupedPhrases []models.TaggedPhrase
		for _, foundIn := range orderedGroups {
			groupedPhrases = append(groupedPhrases, grouped[foundIn]...)
		}

		return groupedPhrases, nil
	}

	return phrases, nil
}

func (p *PhraseModel) Search(ctx context.Context, searchTerm, userID string) ([]models.Phrase, error) {
	searchQuery := "%" + strings.ToLower(searchTerm) + "%"

	query := `
		SELECT id, userId, phrase, phraseDefinition, pinned, foundIn, public, usageCount, createdAt
		FROM phrases
		WHERE userId = ? 
		AND LOWER(phrase || ' ' || phraseDefinition) LIKE ?
	`

	rows, err := p.DB.QueryContext(ctx, query, userID, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("querying search results: %w", err)
	}
	defer rows.Close()

	var phrases []models.Phrase

	for rows.Next() {
		var phrase models.Phrase
		var createdAtStr string

		err = rows.Scan(
			&phrase.ID,
			&phrase.UserID,
			&phrase.Phrase,
			&phrase.PhraseDefinition,
			&phrase.Pinned,
			&phrase.FoundIn,
			&phrase.Public,
			&phrase.UsageCount,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning search row: %w", err)
		}

		phrase.CreatedAt, _ = utils.StringToTime(
			createdAtStr,
			fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase %s", createdAtStr, phrase.Phrase),
		)

		phrases = append(phrases, phrase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating search rows: %w", err)
	}

	return phrases, nil
}

type Scanner interface {
	Scan(dest ...any) error
}

func scanTaggedPhrase(scanner Scanner) (models.Phrase, *models.PhraseTag, error) {
	var phrase models.Phrase
	var phraseCreatedAtStr string
	var tagID, tagPhraseID, tagName, tagColor, tagCreatedAtStr sql.NullString

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
		if tagCreatedAtStr.Valid {
			tag.CreatedAt, _ = utils.StringToTime(
				tagCreatedAtStr.String,
				fmt.Sprintf("Warning: could not parse createdAt '%s' for phrase tag %s for phrase %s", tagCreatedAtStr.String, tag.TagName, phrase.Phrase),
			)
		}
		return phrase, &tag, nil
	}
	return phrase, nil, nil
}
