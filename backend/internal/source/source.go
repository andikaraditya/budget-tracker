package source

import (
	"github.com/andikaraditya/budget-tracker/backend/internal/user"
	"github.com/jackc/pgx/v5/pgtype"
)

type Source struct {
	ID          string             `json:"id"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Initial     float64            `json:"initial"`
	UserId      string             `json:"user_id,omitempty"`
	User        user.UserSimple    `json:"user"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
