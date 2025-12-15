package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JpUnique/goqueue/internal/model"
)

var ErrJobNotFound = errors.New("job not found")

type Postgres struct {
	pool *pgxpool.Pool
}

func (p *Postgres) InsertJob(ctx context.Context, job *model.Job) error {
	sql := `
	INSERT INTO jobs (id, type, payload, status, attempts, max_attempts)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := p.pool.Exec(
		ctx,
		sql,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.Attempts,
		job.MaxAttempts,
	)
	return err
}

func (p *Postgres) GetJob(ctx context.Context, id string) (*model.Job, error) {
	sql := `
	SELECT id, type, payload, status, attempts, max_attempts, created_at, updated_at
	FROM jobs WHERE id = $1
	`
	row := p.pool.QueryRow(ctx, sql, id)

	var j model.Job
	err := row.Scan(
		&j.ID,
		&j.Type,
		&j.Payload,
		&j.Status,
		&j.Attempts,
		&j.MaxAttempts,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return &j, nil
}
