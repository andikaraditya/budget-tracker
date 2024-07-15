package source

import (
	"context"
	"errors"
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

type SourceService interface {
	createSource(req *Source) error
	getSources(userId string, params *apiParams.Params) ([]Source, error)
	getSource(req *Source, userId string) error
	updateSource(req *Source, updatedFields []string) error
}

type srv struct {
	db db.DBService
}

var (
	Service SourceService
)

func init() {
	Service = New(db.Service)
}

func New(db db.DBService) SourceService {
	return &srv{db}
}

func (s *srv) createSource(req *Source) error {
	if req.ID == "" {
		req.ID = uuid.NewString()
	}

	req.CreatedAt = pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
	req.UpdatedAt = req.CreatedAt

	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO "source" (
				id,
				name,
				description,
				initial,
				user_id,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			req.ID,
			req.Name,
			req.Description,
			req.Initial,
			req.UserId,
			req.CreatedAt,
			req.UpdatedAt,
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

	return s.getSource(req, req.UserId)
}

func (s *srv) getSources(userId string, params *apiParams.Params) ([]Source, error) {
	result := []Source{}

	var sb strings.Builder
	args := []any{userId}
	args = params.ComposeFilter(&sb, args)

	rows, err := s.db.Query(
		`SELECT
			id,
			name,
			description,
			initial,
			(SELECT json_build_object(
				'id', id,
				'name', name
				) 
				FROM "user" u
				WHERE u.id = s.user_id
			) AS "user",
			created_at,
			updated_at
		FROM "source" s 
		where user_id = $1 `+sb.String()+params.Sorts.Compose()+params.Page.Compose(),
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Source
		if err := rows.Scan(
			&o.ID,
			&o.Name,
			&o.Description,
			&o.Initial,
			&o.User,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, nil
}
func (s *srv) getSource(req *Source, userId string) error {
	if err := s.db.QueryRow(
		`SELECT
			id,
			name,
			description,
			initial,
			(SELECT json_build_object(
				'id', id,
				'name', name
				) 
				FROM "user" u
				WHERE u.id = s.user_id
			) AS "user",
			created_at,
			updated_at
		FROM "source" s 
		where id = $1 AND user_id = $2;`,
		req.ID,
		userId,
	).Scan(
		&req.ID,
		&req.Name,
		&req.Description,
		&req.Initial,
		&req.User,
		&req.CreatedAt,
		&req.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.ErrNotFound
		}
		return err
	}
	return nil
}
func (s *srv) updateSource(req *Source, updatedFields []string) error {

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
			case "initial":
				args = append(args, req.Initial)
				sb.WriteString(fmt.Sprintf("initial = $%d,", len(args)))
			}
		}

		if len(args) > 2 {
			if _, err := tx.Exec(
				context.Background(),
				fmt.Sprintf(
					`UPDATE source
					SET %s
						updated_at = $2
					WHERE id = $1`,
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
	return s.getSource(req, req.UserId)
}
