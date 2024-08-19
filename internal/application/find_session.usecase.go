package application

import (
	"context"
	"time"
)

type FindSessionUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewFindSessionUseCase(
	sessionRepository SessionRepositoryInterface,
) *FindSessionUseCase {
	return &FindSessionUseCase{
		SessionRepository: sessionRepository,
	}
}

type FindSessionInputDto struct {
	ID string `json:"id"`
}

type FindSessionOutputDto struct {
	ID               string        `json:"id"`
	Price            float64       `json:"price"`
	Notes            string        `json:"notes"`
	Date             time.Time     `json:"date"`
	Duration         time.Duration `json:"duration"`
	PatientName      string        `json:"patient_name"`
	PatientId        string        `json:"patient_id"`
	Plan             string        `json:"plan"`
	ProfessionalName string        `json:"professional_name"`
	ProfessionalId   string        `json:"professional_id"`
}

func (u *FindSessionUseCase) Execute(ctx context.Context, input FindSessionInputDto) (FindSessionOutputDto, error) {
	session, err := u.SessionRepository.Find(ctx, FindSessionRepositoryInput{
		ID: input.ID,
	})

	dto := FindSessionOutputDto{}
	if err != nil {
		return dto, err
	}

	dto.ID = session.ID
	dto.Price = session.Price
	dto.Notes = session.Notes
	dto.Date = session.Date
	dto.Duration = session.Duration
	dto.PatientName = session.Patient.Name
	dto.ProfessionalName = session.Professional.Name
	dto.Plan = session.Plan
	dto.PatientId = session.Patient.ID
	dto.ProfessionalId = session.Professional.ID
	return dto, nil
}
