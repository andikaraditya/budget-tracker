package transfer

import (
	"github.com/andikaraditya/budget-tracker/backend/internal/source"
	"github.com/andikaraditya/budget-tracker/backend/internal/user"
	"github.com/jackc/pgx/v5/pgtype"
)

type Transfer struct {
	ID        string             `json:"id"`
	ToID      string             `json:"to_id,omitempty" validate:"required"`
	To        source.Source      `json:"to" validate:"-"`
	FromId    string             `json:"from_id,omitempty" validate:"required"`
	From      source.Source      `json:"from" validate:"-"`
	Amount    float64            `json:"amount" validate:"required"`
	UserId    string             `json:"user_id,omitempty"`
	User      user.UserSimple    `json:"user"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
