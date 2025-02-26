package models

import (
	"github.com/arinji2/vocab-thing/internal/tools/types"
)

type User struct {
	ID        int            `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	CreatedAt types.DateTime `json:"createdAt"`
}
