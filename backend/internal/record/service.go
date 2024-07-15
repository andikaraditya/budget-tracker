package record

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	"github.com/andikaraditya/budget-tracker/backend/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RecordService interface {
	createRecord(req *Record) error
	getRecords(userId string) ([]Record, error)
	getRecord(req *Record) error
	updateRecord(req *Record, updatedFields []string) error
	getSummary(req *Summary, userId string) error
}

type srv struct {
	db db.DBService
}

var (
	Service RecordService
)

func init() {
	Service = New(db.Service)
}

func New(db db.DBService) RecordService {
	return &srv{db}
}

func (s *srv) createRecord(req *Record) error {
	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		req.ID = uuid.NewString()
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO "record" (
				id,
				amount,
				description,
				type,
				category_id,
				source_id,
				user_id,
				date
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			req.ID,
			req.Amount,
			req.Description,
			req.Type,
			req.CategoryId,
			req.SourceId,
			req.UserId,
			req.Date,
		)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return s.getRecord(req)
}

func (s *srv) getRecords(userId string) ([]Record, error) {
	var r []Record
	rows, err := s.db.Query(
		`SELECT 
			id,
			amount,
			description,
			"type" ,
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'type', type,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "category" a
					WHERE a.id = r.category_id
			) AS "category",
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "source" b
					WHERE b.id = r.source_id
			) AS "source",
			(SELECT json_build_object(
					'id', id,
					'name', name
					) 
					FROM "user" u
					WHERE u.id = r.user_id
			) AS "user",
			"date",
			created_at,
			updated_at
		FROM record r
		WHERE r.user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var o Record
		if err := rows.Scan(
			&o.ID,
			&o.Amount,
			&o.Description,
			&o.Type,
			&o.Category,
			&o.Source,
			&o.User,
			&o.Date,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		r = append(r, o)
	}
	return r, nil
}
func (s *srv) getRecord(req *Record) error {
	if err := s.db.QueryRow(
		`SELECT 
			id,
			amount,
			description,
			"type" ,
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'type', type,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "category" a
					WHERE a.id = r.category_id
			) AS "category",
			(SELECT json_build_object(
					'id', id,
					'name', name,
					'description', description,
					'created_at', created_at,
					'updated_at', updated_at
					) 
					FROM "source" b
					WHERE b.id = r.source_id
			) AS "source",
			(SELECT json_build_object(
					'id', id,
					'name', name
					) 
					FROM "user" u
					WHERE u.id = r.user_id
			) AS "user",
			"date",
			created_at,
			updated_at
		FROM record r
		WHERE r.id = $1 AND r.user_id = $2`,
		req.ID,
		req.UserId,
	).Scan(
		&req.ID,
		&req.Amount,
		&req.Description,
		&req.Type,
		&req.Category,
		&req.Source,
		&req.User,
		&req.Date,
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

func (s *srv) updateRecord(req *Record, updatedFields []string) error {
	if err := s.db.Commit(nil, func(tx pgx.Tx) error {
		req.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		args := []any{req.ID, req.UpdatedAt}
		var sb strings.Builder

		for _, field := range updatedFields {
			switch field {
			case "amount":
				args = append(args, req.Amount)
				sb.WriteString(fmt.Sprintf("amount = $%d,", len(args)))
			case "description":
				args = append(args, req.Description)
				sb.WriteString(fmt.Sprintf("description = $%d,", len(args)))
			case "date":
				args = append(args, req.Date)
				sb.WriteString(fmt.Sprintf("date = $%d,", len(args)))
			}
		}

		if len(args) > 2 {
			if _, err := tx.Exec(
				context.Background(),
				fmt.Sprintf(
					`UPDATE record
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
	return s.getRecord(req)
}

func (s *srv) getSummary(req *Summary, userId string) error {
	if err := s.db.QueryRow(
		`WITH totals AS (
			SELECT 
				SUM(CASE WHEN r."type" = 'expense' THEN r.amount ELSE 0 END) AS expense,
				SUM(CASE WHEN r."type" = 'income' THEN r.amount ELSE 0 END) AS income
			FROM record r
			WHERE r.user_id = $1
		)
		SELECT 
			expense,
			income,
			(income - expense) AS total
		FROM totals`,
		userId,
	).Scan(
		&req.Expense,
		&req.Income,
		&req.Total,
	); err != nil {
		return err
	}
	return nil
}
