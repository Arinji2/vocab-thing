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

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, createdAt FROM users;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Printf("querying users: %s", err.Error())
		return nil, errorcode.ErrDBQuery
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var createdAtStr string

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
		if err != nil {
			log.Printf("scanning user row: %s", err.Error())
			return nil, errorcode.ErrScanningRow
		}

		parsedTime, err := utils.StringToTime(createdAtStr)
		if err != nil {
			log.Printf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID)
			parsedTime = time.Now().UTC()
		}
		user.CreatedAt = parsedTime

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterating user rows: %s", err.Error())
		return nil, errorcode.ErrIteratingRows
	}

	return users, nil
}

func (m *UserModel) ByID(ctx context.Context, id string) (models.User, error) {
	query := `SELECT id, username, email, createdAt FROM users WHERE id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User
	var createdAtStr string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("user not found with id %s: %s", id, err.Error())
			return models.User{}, errorcode.ErrDBQuery
		}
		log.Printf("scanning user row: %s", err.Error())
		return models.User{}, errorcode.ErrScanningRow
	}

	parsedTime, err := utils.StringToTime(createdAtStr)
	if err != nil {
		log.Printf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID)
		parsedTime = time.Now().UTC()
	}

	user.CreatedAt = parsedTime

	return user, nil
}

func (m *UserModel) ByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT id, username, email, createdAt FROM users WHERE username = ?`

	row := m.DB.QueryRowContext(ctx, query, username)

	var user models.User
	var createdAtStr string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("user not found with username %s: %s", username, err.Error())
			return models.User{}, errorcode.ErrDBQuery
		}
		log.Printf("scanning user row: %s", err.Error())
		return models.User{}, errorcode.ErrScanningRow
	}

	parsedTime, err := utils.StringToTime(createdAtStr)
	if err != nil {
		log.Printf("could not parse createdAt '%s' for user %s: %s", createdAtStr, user.ID, err.Error())
		parsedTime = time.Now().UTC()
	}
	user.CreatedAt = parsedTime

	return user, nil
}

func (m *UserModel) ByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT id, username, email, createdAt FROM users WHERE email = ?`

	row := m.DB.QueryRowContext(ctx, query, email)

	var user models.User
	var createdAtStr string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("user not found with email %s: %s", email, err.Error())
			return models.User{}, errorcode.ErrDBQuery
		}
		log.Printf("scanning user row: %s", err.Error())
		return models.User{}, errorcode.ErrScanningRow
	}

	parsedTime, err := utils.StringToTime(createdAtStr)
	if err != nil {
		log.Printf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID)
		parsedTime = time.Now().UTC()
	}
	user.CreatedAt = parsedTime

	return user, nil
}

func (m *UserModel) Create(ctx context.Context, user *models.User) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("starting transaction: %s", err.Error())
		return errorcode.ErrTransactionStart
	}
	defer tx.Rollback()

	user.CreatedAt = time.Now().UTC()
	query := `INSERT INTO users (id, username, email, createdAt) 
          VALUES (lower(hex(randomblob(16))), ?, ?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, user.Username, user.Email, user.CreatedAt.Format(time.RFC3339)).Scan(&user.ID)
	if err != nil {
		log.Printf("error with user creation of username %s: %s", user.Username, err.Error())
		return errorcode.ErrDBCreate
	}

	if err := tx.Commit(); err != nil {
		log.Printf("committing transaction: %s", err.Error())
		return errorcode.ErrTransactionCommit
	}

	return nil
}
