package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/repository"
)

type PatientService struct {
	patientRepository *repository.PatientRepository
}

func NewPatientService(
	patientRepository *repository.PatientRepository,
) *PatientService {
	return &PatientService{
		patientRepository: patientRepository,
	}
}

func (s *PatientService) Create(
	ctx context.Context,
	patient model.Patient,
) (model.PatientResponseBody, error) {
	userIdString := ctx.Value("userID").(string)
	userId, err := uuid.Parse(
		userIdString,
	)
	if err != nil {
		return model.PatientResponseBody{}, constant.ErrUnauthorized
	}

	currentTime := time.Now()
	patient.CreatedAt = currentTime
	patient.UpdatedAt = currentTime
	patient.UserID = userId
	saved, err := s.patientRepository.Create(
		ctx,
		patient,
	)
	if err != nil {
		return model.PatientResponseBody{}, err
	}

	return saved.ToResponseBody()
}

func (s *PatientService) FindAll(
	ctx context.Context,
	queries model.PatientQuery,
) ([]model.PatientResponseBody, error) {
	patients, err := s.patientRepository.FindAll(
		ctx,
		queries,
	)
	if err != nil {
		return nil, err
	}

	patientData := make(
		[]model.PatientResponseBody,
		0,
		len(patients),
	)
	for _, patient := range patients {
		data, err := patient.ToResponseBody()
		if err != nil {
			return nil, err
		}

		patientData = append(
			patientData,
			data,
		)
	}

	return patientData, nil
}
