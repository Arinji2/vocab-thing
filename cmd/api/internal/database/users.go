package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var createdAtStr string

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("scanning user row: %w", err)
		}
		parsedTime, err := utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID))
		if err != nil {
			fmt.Println(err.Error())
		}
		user.CreatedAt = parsedTime

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user rows: %w", err)
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
			return models.User{}, fmt.Errorf("user not found with id %s: %w", id, err)
		}
		return models.User{}, fmt.Errorf("scanning user row: %w", err)
	}

	parsedTime, err := utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID))
	if err != nil {
		fmt.Println(err.Error())
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
			return models.User{}, fmt.Errorf("user not found with username %s: %w", username, err)
		}
		return models.User{}, fmt.Errorf("scanning user row: %w", err)
	}

	parsedTime, err := utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID))
	if err != nil {
		fmt.Println(err.Error())
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
			return models.User{}, fmt.Errorf("user not found with email %s: %w", email, err)
		}
		return models.User{}, fmt.Errorf("scanning user row: %w", err)
	}

	parsedTime, err := utils.StringToTime(createdAtStr, fmt.Sprintf("Warning: could not parse createdAt '%s' for user %s", createdAtStr, user.ID))
	if err != nil {
		fmt.Println(err.Error())
	}
	user.CreatedAt = parsedTime

	return user, nil
}

func (m *UserModel) Create(ctx context.Context, user models.User) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	user.CreatedAt = time.Now().UTC()
	query := `INSERT INTO users (id, username, email, createdAt) 
          VALUES (lower(hex(randomblob(16))), ?, ?, ?)`

	_, err = tx.ExecContext(ctx, query, user.Username, user.Email, user.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error with user creation of username %s, userID %s: %w", user.Username, user.ID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
