package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/JpUnique/goqueue/internal/model"
	"github.com/JpUnique/goqueue/internal/store"
)

type JobService struct {
	pg *store.Postgres
	rd *store.RedisStore
}

func NewJobService(pg *store.Postgres, rd *store.RedisStore) *JobService {
	return &JobService{pg: pg, rd: rd}
}

func (s *JobService) CreateJob(
	ctx context.Context,
	jobType string,
	payload any,
) (string, error) {

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	job := &model.Job{
		ID:          uuid.NewString(),
		Type:        jobType,
		Payload:     data,
		Status:      model.StatusPending,
		Attempts:    0,
		MaxAttempts: 3,
	}

	// 1️Persist first
	if err := s.pg.InsertJob(ctx, job); err != nil {
		return "", err
	}

	// 2️Enqueue second
	if err := s.rd.Enqueue(ctx, job.ID); err != nil {
		return "", err
	}

	return job.ID, nil
}

func (s *JobService) GetJob(ctx context.Context, id string) (*model.Job, error) {
	return s.pg.GetJob(ctx, id)
}
