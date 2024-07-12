package db

import (
	"context"
	"fmt"

	logrusAdapter "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
)

const (
	OrderAsc  = "ASC"
	OrderDesc = "DESC NULLS LAST"

	UniqueViolation = "23505"
)

type txFn func(pgx.Tx) error

type DBService interface {
	Query(q string, args ...any) (pgx.Rows, error)
	QueryRow(q string, args ...any) pgx.Row
	Commit(tx pgx.Tx, fn txFn) error
}

type impl struct {
	db *pgxpool.Pool
}

var Service DBService

func init() {
	Service = New()
}

func New() DBService {
	pgxCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	pgxCfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   logrusAdapter.NewLogger(logrus.New()),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	return &impl{db: pool}
}

func (s *impl) Query(q string, args ...any) (pgx.Rows, error) {
	return s.db.Query(context.Background(), q, args...)
}

func (s *impl) QueryRow(q string, args ...any) pgx.Row {
	return s.db.QueryRow(context.Background(), q, args...)
}

func (s *impl) Commit(tx pgx.Tx, fn txFn) error {
	given := tx != nil

	var err error

	if !given {
		tx, err = s.db.Begin(context.Background())
		if err != nil {
			fmt.Printf("Error: %v", err)
			return err
		}
	}

	defer func() {
		if !given && err != nil {
			if err := tx.Rollback(context.Background()); err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
	}()

	if err = fn(tx); err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	if !given {
		if err = tx.Commit(context.Background()); err != nil {
			fmt.Printf("Error: %v", err)

			return err
		}
	}
	return nil
}
