package category

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	"github.com/andikaraditya/budget-tracker/backend/internal/db"
	apiParams "github.com/andikaraditya/budget-tracker/backend/internal/params"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CategoryService interface {
	createCategory(req *Category) error
	getCategories(params *apiParams.Params) ([]Category, error)
	getCategory(req *Category) error
	updateCategory(req *Category, updatedFields []string) error
}

type srv struct {
	db db.DBService
}

var (
	Service CategoryService
)

func init() {
	Service = New(db.Service)
}

func New(db db.DBService) CategoryService {
	return &srv{db}
}

func (s *srv) createCategory(req *Category) error {
	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		req.ID = uuid.NewString()
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO "category" (
				id,
				name,
				description,
				type
			) VALUES ($1, $2, $3, $4);`,
			req.ID,
			req.Name,
			req.Description,
			req.Type,
		)
		if err != nil {
			if strings.Contains(err.Error(), "23505") {
				return api.ErrPayload
			}
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return s.getCategory(req)
}

func (s *srv) getCategories(params *apiParams.Params) ([]Category, error) {
	var c []Category
	var sb strings.Builder
	args := []any{}
	args = params.ComposeFilter(&sb, args)

	rows, err := s.db.Query(
		`SELECT 
			id,
			name,
			description,
			type,
			created_at,
			updated_at
		FROM category
		WHERE TRUE `+sb.String()+params.Sorts.Compose()+params.Page.Compose(),
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Category
		if err := rows.Scan(
			&o.ID,
			&o.Name,
			&o.Description,
			&o.Type,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}

		c = append(c, o)
	}
	return c, nil
}
func (s *srv) getCategory(req *Category) error {
	s.db.QueryRow(
		`SELECT 
			id,
			name,
			description,
			type,
			created_at,
			updated_at
		FROM category
		WHERE id = $1;`,
		req.ID,
	).Scan(
		&req.ID,
		&req.Name,
		&req.Description,
		&req.Type,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	return nil
}
func (s *srv) updateCategory(req *Category, updatedFields []string) error {
	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		req.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		args := []any{req.ID, req.UpdatedAt}
		var sb strings.Builder

		for _, field := range updatedFields {
			switch field {
			case "name":
				args = append(args, req.Name)
				sb.WriteString(fmt.Sprintf("name = $%d,", len(args)))
			case "description":
				args = append(args, req.Description)
				sb.WriteString(fmt.Sprintf("description = $%d,", len(args)))
			}
		}

		if len(args) > 2 {
			if _, err := tx.Exec(
				context.Background(),
				fmt.Sprintf(
					`UPDATE category
					SET %s
						updated_at = $2
					WHERE id = $1;`,
					sb.String(),
				),
				args...,
			); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return s.getCategory(req)
}
