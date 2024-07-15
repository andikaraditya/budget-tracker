package transfer

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

type TransferService interface {
	createTransfer(req *Transfer) error
	getTransfers(userId string, params *apiParams.Params) ([]Transfer, error)
	getTransfer(req *Transfer, userId string) error
	updateTransfer(req *Transfer, updatedFields []string) error
}

type srv struct {
	db db.DBService
}

var (
	Service TransferService
)

func init() {
	Service = New(db.Service)
}

func New(db db.DBService) TransferService {
	return &srv{db}
}

func (s *srv) createTransfer(req *Transfer) error {
	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		req.ID = uuid.NewString()
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO "transfer" (
				id,
				to_id,
				from_id,
				amount,
				user_id
			) VALUES ($1, $2, $3, $4, $5)`,
			req.ID,
			req.ToID,
			req.FromId,
			req.Amount,
			req.UserId,
		)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return s.getTransfer(req, req.UserId)
}

func (s *srv) getTransfers(userId string, params *apiParams.Params) ([]Transfer, error) {
	var t []Transfer

	var sb strings.Builder
	args := []any{userId}
	args = params.ComposeFilter(&sb, args)

	rows, err := s.db.Query(
		`SELECT 
			id,
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "source" a
					WHERE a.id = t.to_id
			) AS "to",
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "source" b
					WHERE b.id = t.from_id
			) AS "from",
			amount,
			(SELECT json_build_object(
					'id', id,
					'name', name
					) 
					FROM "user" u
					WHERE u.id = t.user_id
			) AS "user",
				created_at,
				updated_at
		FROM transfer t 
		WHERE t.user_id = $1 `+sb.String()+params.Sorts.Compose()+params.Page.Compose(),
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Transfer
		rows.Scan(
			&o.ID,
			&o.To,
			&o.From,
			&o.Amount,
			&o.User,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		t = append(t, o)
	}
	return t, nil
}

func (s *srv) getTransfer(req *Transfer, userId string) error {
	if err := s.db.QueryRow(
		`SELECT 
			id,
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "source" a
					WHERE a.id = t.to_id
			) AS "to",
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "source" b
					WHERE b.id = t.from_id
			) AS "from",
			amount,
			(SELECT json_build_object(
					'id', id,
					'name', name
					) 
					FROM "user" u
					WHERE u.id = t.user_id
			) AS "user",
				created_at,
				updated_at
		FROM transfer t 
		WHERE t.id = $1 AND t.user_id = $2`,
		req.ID,
		userId,
	).Scan(
		&req.ID,
		&req.To,
		&req.From,
		&req.Amount,
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

func (s *srv) updateTransfer(req *Transfer, updatedFields []string) error {
	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		req.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		args := []any{req.ID, req.UpdatedAt}
		var sb strings.Builder

		for _, field := range updatedFields {
			switch field {
			case "to_id":
				args = append(args, req.ToID)
				sb.WriteString(fmt.Sprintf("to_id = $%d,", len(args)))
			case "from_id":
				args = append(args, req.FromId)
				sb.WriteString(fmt.Sprintf("from_id = $%d,", len(args)))
			case "amount":
				args = append(args, req.Amount)
				sb.WriteString(fmt.Sprintf("amount = $%d,", len(args)))
			}
		}
		_, err := tx.Exec(
			context.Background(),
			fmt.Sprintf(
				`UPDATE transfer
				SET %s
					updated_at = $2
				WHERE id = $1`,
				sb.String(),
			),
			args...,
		)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return s.getTransfer(req, req.UserId)
}
