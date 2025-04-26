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
        SELECT p.id, p.userId, p.phrase, p.phraseDefinition, p.pinned, p.foundIn, p.public, p.usageCount, p.createdAt, p.updatedAt, p.deletedAt,
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
	if taggedPhrase.Tag == nil {
		taggedPhrase.Tag = []models.PhraseTag{}
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
	if sortBy == "" {
		sortBy = "createdAt"
	}
	if order == "" || (strings.ToUpper(order) != "ASC" && strings.ToUpper(order) != "DESC") {
		order = "DESC"
	}

	totalPages, err := p.CountTotalPages(ctx, pageSize, userID)
	if err != nil {
		log.Printf("Error counting total pages for user %s: %v", userID, err)
		return nil, errorcode.ErrDBQuery
	}

	if totalPages == 0 {
		return []models.TaggedPhrase{}, nil
	}

	if pageNumber > totalPages {
		return nil, fmt.Errorf("page number %d is greater than total pages %d", pageNumber, totalPages)
	}

	query := fmt.Sprintf(`
		SELECT
			p.id, p.userId, p.phrase, p.phraseDefinition, p.pinned, p.foundIn, p.public, p.usageCount, p.createdAt, p.updatedAt, p.deletedAt,
			pt.id, pt.phraseId, pt.tagName, pt.tagColor, pt.createdAt
		FROM phrases p
		LEFT JOIN phrase_tags pt ON p.id = pt.phraseId
		WHERE p.userId = ? AND p.deletedAt IS NULL
		ORDER BY p.%s %s, p.id
		LIMIT ? OFFSET ?
	`, sortBy, order)

	offset := (pageNumber - 1) * pageSize

	rows, err := p.DB.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		log.Printf("Error querying paginated tagged phrases for user %s: %v", userID, err)
		return nil, errorcode.ErrDBQuery
	}
	defer rows.Close()

	taggedPhrasesMap := make(map[string]*models.TaggedPhrase)
	orderedPhraseIDs := []string{}

	for rows.Next() {
		phrase, tag, err := scanTaggedPhrase(rows)
		if err != nil {
			log.Printf("Error scanning tagged phrase row during All() for user %s: %v", userID, err)
			continue
		}

		if _, exists := taggedPhrasesMap[phrase.ID]; !exists {
			taggedPhrasesMap[phrase.ID] = &models.TaggedPhrase{
				Phrase: phrase,
				Tag:    []models.PhraseTag{},
			}
			orderedPhraseIDs = append(orderedPhraseIDs, phrase.ID)
		}

		if tag != nil {
			if taggedPhrasesMap[phrase.ID].Tag == nil {
				taggedPhrasesMap[phrase.ID].Tag = []models.PhraseTag{}
			}
			taggedPhrasesMap[phrase.ID].Tag = append(taggedPhrasesMap[phrase.ID].Tag, *tag)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating tagged phrase rows during All() for user %s: %v", userID, err)
		return nil, errorcode.ErrDBQuery
	}

	finalPhrases := make([]models.TaggedPhrase, 0, len(orderedPhraseIDs))
	for _, id := range orderedPhraseIDs {
		finalPhrases = append(finalPhrases, *taggedPhrasesMap[id])
	}

	if groupBy == "foundIn" || groupBy == "public" {
		grouped := make(map[string][]models.TaggedPhrase)
		orderedGroups := []string{}

		for _, taggedPhrase := range finalPhrases {
			var key string
			switch groupBy {
			case "foundIn":
				key = taggedPhrase.Phrase.FoundIn
				if key == "" {
					key = "[Unknown Source]"
				}
			case "public":
				key = strconv.FormatBool(taggedPhrase.Phrase.Public)
			}

			if _, exists := grouped[key]; !exists {
				orderedGroups = append(orderedGroups, key)
			}
			grouped[key] = append(grouped[key], taggedPhrase)
		}

		var groupedPhrases []models.TaggedPhrase
		for _, groupKey := range orderedGroups {
			groupedPhrases = append(groupedPhrases, grouped[groupKey]...)
		}
		return groupedPhrases, nil
	}

	return finalPhrases, nil
}

func (p *PhraseModel) Search(ctx context.Context, searchTerm, userID string) ([]models.Phrase, error) {
	searchQuery := "%" + strings.ToLower(searchTerm) + "%"

	query := `
		SELECT id, userId, phrase, phraseDefinition, pinned, foundIn, public, usageCount, createdAt, updatedAt, deletedAt
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
            UPDATE phrases SET phrase = ?, phraseDefinition = ?, pinned = ?, foundIn = ?, public = ?, usageCount = ?, updatedAt = ?
            WHERE id = ? AND userId = ?
            `

	res, err := tx.ExecContext(ctx, query,
		phrase.Phrase, phrase.PhraseDefinition, phrase.Pinned,
		phrase.FoundIn, phrase.Public, phrase.UsageCount, time.Now().UTC(),
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
            UPDATE phrases SET deletedAt =?
            WHERE id = ? AND userId = ?
            `

	res, err := tx.ExecContext(ctx, query, time.Now().UTC(), phraseID, userID)
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
	var phraseCreatedAtStr, phraseUpdatedAtStr string
	var phraseDeletedAtStr, tagID, tagPhraseID, tagName, tagColor, tagCreatedAtStr sql.NullString

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
		&phraseUpdatedAtStr,
		&phraseDeletedAtStr,
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

	phrase.UpdatedAt, err = utils.StringToTime(phraseUpdatedAtStr)
	if err != nil {
		log.Printf("Warning: could not parse updatedAt '%s' for phrase %s", phraseUpdatedAtStr, phrase.Phrase)
		phrase.UpdatedAt = time.Now().UTC()
	}

	if phraseDeletedAtStr.Valid {
		deletedAtTime, err := utils.StringToTime(phraseDeletedAtStr.String)
		if err != nil {
			log.Printf("Warning: could not parse deletedAt '%s' for phrase %s", phraseDeletedAtStr.String, phrase.Phrase)
			phrase.DeletedAt = nil
		}
		phrase.DeletedAt = &deletedAtTime
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
