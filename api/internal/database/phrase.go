package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/arinji2/vocab-thing/internal/errorcode"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils"
)

type PhraseModel struct {
	DB *sql.DB
}

func (p *PhraseModel) CreatePhrase(ctx context.Context, phrase *models.Phrase) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return errorcode.ErrTransactionStart
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
		log.Printf("error with phrase creation of userID %s: %s", phrase.UserID, err.Error())
		return errorcode.ErrPhraseCreation
	}
	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return errorcode.ErrTransactionCommit
	}

	return nil
}

func (p *PhraseModel) CreateTag(ctx context.Context, phrase *models.PhraseTag) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return errorcode.ErrTransactionStart
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
		log.Printf("error with phrase tag creation of phraseID %s: %s", phrase.PhraseID, err.Error())
		return errorcode.ErrPhraseTagCreation
	}
	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return errorcode.ErrTransactionCommit
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
		log.Printf("querying phrase: %s", err.Error())
		return nil, errorcode.ErrDBQuery
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
	if totalPages == 0 {
		return []models.TaggedPhrase{}, nil
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

		phrase.CreatedAt, err = utils.StringToTime(phraseCreatedAtStr)
		if err != nil {
			log.Printf("Warning: could not parse createdAt '%s' for phrase %s", phraseCreatedAtStr, phrase.Phrase)
			phrase.CreatedAt = time.Now().UTC()
		}

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

		tag.CreatedAt, err = utils.StringToTime(tagCreatedAtStr)
		if err != nil {
			log.Printf("Warning: could not parse createdAt '%s' for phrase tag %s", tagCreatedAtStr, tag.TagName)
			tag.CreatedAt = time.Now().UTC()
		}

		if idx, exists := phraseIndex[tag.PhraseID]; exists {
			phrases[idx].Tag = append(phrases[idx].Tag, tag)
		}
	}

	if err := tagRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tag rows: %w", err)
	}
	if groupBy == "foundIn" || groupBy == "public" {
		grouped := make(map[string][]models.TaggedPhrase)
		orderedGroups := []string{}

		for _, phrase := range phrases {
			switch groupBy {
			case "foundIn":
				if _, exists := grouped[phrase.Phrase.FoundIn]; !exists {
					orderedGroups = append(orderedGroups, phrase.Phrase.FoundIn)
				}
				grouped[phrase.Phrase.FoundIn] = append(grouped[phrase.Phrase.FoundIn], phrase)

			case "public":
				if _, exists := grouped[strconv.FormatBool(phrase.Phrase.Public)]; !exists {
					orderedGroups = append(orderedGroups, strconv.FormatBool(phrase.Phrase.Public))
				}
				grouped[strconv.FormatBool(phrase.Phrase.Public)] = append(grouped[strconv.FormatBool(phrase.Phrase.Public)], phrase)
			}
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

		phrase.CreatedAt, err = utils.StringToTime(createdAtStr)
		if err != nil {
			log.Printf("Warning: could not parse createdAt '%s' for phrase %s", createdAtStr, phrase.Phrase)
			phrase.CreatedAt = time.Now().UTC()
		}

		phrases = append(phrases, phrase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating search rows: %w", err)
	}

	if len(phrases) == 0 {
		return []models.Phrase{}, nil
	}
	return phrases, nil
}

func (p *PhraseModel) UpdatePhrase(ctx context.Context, phrase *models.Phrase, userID string) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer tx.Rollback()
	phrase.CreatedAt = time.Now().UTC()
	query := `
            UPDATE phrases SET phrase = ?, phraseDefinition = ?, pinned = ?, foundIn = ?, public = ?, usageCount = ?
            WHERE id = ? AND userId = ?
            `

	res, err := tx.ExecContext(ctx, query,
		phrase.Phrase, phrase.PhraseDefinition, phrase.Pinned,
		phrase.FoundIn, phrase.Public, phrase.UsageCount,
		phrase.ID, userID,
	)
	if err != nil {
		return fmt.Errorf("error updating phrase (id: %s, userID: %s): %w", phrase.ID, phrase.UserID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no phrase found with id: %s for user: %s", phrase.ID, userID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *PhraseModel) UpdateTag(ctx context.Context, tag *models.PhraseTag, userID string) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()
	tag.CreatedAt = time.Now().UTC()
	query := `
    UPDATE phrase_tags 
    SET tagName = ?, tagColor = ?
    WHERE id = ? AND phraseId = ?
    AND EXISTS (
      SELECT 1 FROM phrases 
      WHERE phrases.id = ? 
      AND phrases.userId = ?
    );
                `

	res, err := tx.ExecContext(ctx, query,
		tag.TagName, tag.TagColor, tag.ID, tag.PhraseID,
		tag.PhraseID, userID,
	)
	if err != nil {
		return fmt.Errorf("error updating tag (id: %s, phraseId: %s, userId: %s): %w", tag.ID, tag.PhraseID, userID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no tag found with id: %s for user: %s", tag.ID, userID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *PhraseModel) DeletePhrase(ctx context.Context, phraseID, userID string) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer tx.Rollback()
	query := `
            DELETE FROM phrases WHERE id = ? AND userId = ?
            `

	res, err := tx.ExecContext(ctx, query, phraseID, userID)
	if err != nil {
		return fmt.Errorf("error deleting phrase (id: %s, userID: %s): %w", phraseID, userID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no phrase deleted with id: %s for user: %s", phraseID, userID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *PhraseModel) DeleteTag(ctx context.Context, phraseID, tagID, userID string) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer tx.Rollback()
	query := `
DELETE FROM phrase_tags 
		WHERE id = ? 
    AND phraseId = ? 
		AND EXISTS (
			SELECT 1 FROM phrases 
			WHERE phrases.id = phrase_tags.phraseId 
			AND phrases.id = ? 
			AND phrases.userId = ?
		);
	`

	res, err := tx.ExecContext(ctx, query, tagID, phraseID, phraseID, userID)
	if err != nil {
		return fmt.Errorf("error deleting tag (id: %s, phraseID: %s,  userID: %s): %w", tagID, phraseID, userID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no tag deleted with id: %s for user: %s", tagID, userID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTaggedPhrase(scanner scanner) (models.Phrase, *models.PhraseTag, error) {
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

	phrase.CreatedAt, err = utils.StringToTime(phraseCreatedAtStr)
	if err != nil {
		log.Printf("Warning: could not parse createdAt '%s' for phrase %s", phraseCreatedAtStr, phrase.Phrase)
		phrase.CreatedAt = time.Now().UTC()
	}

	if tagID.Valid && tagName.Valid {
		tag := models.PhraseTag{
			ID:       tagID.String,
			PhraseID: tagPhraseID.String,
			TagName:  tagName.String,
			TagColor: tagColor.String,
		}
		if tagCreatedAtStr.Valid {
			tag.CreatedAt, err = utils.StringToTime(tagCreatedAtStr.String)
			if err != nil {
				log.Printf("Warning: could not parse createdAt '%s' for phrase tag %s for phrase %s", tagCreatedAtStr.String, tag.TagName, phrase.Phrase)
				tag.CreatedAt = time.Now().UTC()
			}

		}
		return phrase, &tag, nil
	}
	return phrase, nil, nil
}
