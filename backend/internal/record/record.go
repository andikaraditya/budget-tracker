package record

import (
	"github.com/andikaraditya/budget-tracker/backend/internal/category"
	"github.com/andikaraditya/budget-tracker/backend/internal/source"
	"github.com/andikaraditya/budget-tracker/backend/internal/user"
	"github.com/jackc/pgx/v5/pgtype"
)

type Record struct {
	ID          string             `json:"id"`
	Amount      float64            `json:"amount" validate:"required"`
	Description string             `json:"description"`
	Type        string             `json:"type" validate:"required,oneof=income expense"`
	CategoryId  string             `json:"category_id,omitempty" validate:"required"`
	Category    category.Category  `json:"category,omitempty" validate:"-"`
	SourceId    string             `json:"source_id,omitempty" validate:"required"`
	Source      source.Source      `json:"source" validate:"-"`
	UserId      string             `json:"-"`
	User        user.UserSimple    `json:"user"`
	Date        pgtype.Timestamptz `json:"date" validate:"required"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type Summary struct {
	Expense float64 `json:"expense"`
	Income  float64 `json:"income"`
	Total   float64 `json:"total"`
}
