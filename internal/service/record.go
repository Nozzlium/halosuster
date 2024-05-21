package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/repository"
)

type RecordService struct {
	recordRepository *repository.RecordRepository
}

func NewRecordService(
	recordRepository *repository.RecordRepository,
) *RecordService {
	return &RecordService{
		recordRepository: recordRepository,
	}
}

func (s *RecordService) Create(
	ctx context.Context,
	record model.Record,
) (model.Record, error) {
	userIdString := ctx.Value("userID").(string)
	userId, err := uuid.Parse(
		userIdString,
	)
	if err != nil {
		return model.Record{}, constant.ErrUnauthorized
	}

	currentTime := time.Now()
	record.CreatedAt = currentTime
	record.UpdatedAt = currentTime

	record.UserID = userId
	id, err := uuid.NewV7()
	if err != nil {
		return model.Record{}, err
	}
	record.ID = id
	saved, err := s.recordRepository.Create(
		ctx,
		record,
	)
	if err != nil {
		return model.Record{}, err
	}

	return saved, nil
}
