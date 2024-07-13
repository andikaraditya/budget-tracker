package category

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Category struct {
	ID          string             `json:"id"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Type        string             `json:"type" validate:"required,oneof=income expense"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
