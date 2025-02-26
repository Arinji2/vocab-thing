package models

import "github.com/arinji2/vocab-thing/internal/utils/datetime"

type User struct {
	ID        int               `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	CreatedAt datetime.DateTime `json:"createdAt"`
}
