package application

import (
	"context"
	"time"
)

type UpdateSessionUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewUpdateSessionUseCase(
	sessionRepository SessionRepositoryInterface,
) *UpdateSessionUseCase {
	return &UpdateSessionUseCase{
		SessionRepository: sessionRepository,
	}
}

type UpdateSessionInputDTO struct {
	Notes string `json:"notes"`
	ID    string `json:"id"`
}

type UpdateSessionOutputDTO struct {
	ID               string        `json:"id"`
	Price            float64       `json:"price"`
	Notes            string        `json:"notes"`
	Date             time.Time     `json:"date"`
	Duration         time.Duration `json:"duration"`
	PatientName      string        `json:"patient_name"`
	ProfessionalName string        `json:"professional_name"`
}

func (u *UpdateSessionUseCase) Execute(ctx context.Context, input UpdateSessionInputDTO) (UpdateSessionOutputDTO, error) {
	dto := UpdateSessionOutputDTO{}
	repoInput := FindSessionRepositoryInput{
		ID: input.ID,
	}
	session, err := u.SessionRepository.Find(ctx, repoInput)
	if err != nil {
		return dto, err
	}
	err = session.ChangeNote(input.Notes)
	if err != nil {
		return dto, err
	}
	err = u.SessionRepository.Update(ctx, session)
	if err != nil {
		return dto, err
	}
	dto.ID = session.ID
	dto.Price = session.Price
	dto.Notes = session.Notes
	dto.Date = session.Date
	dto.Duration = session.Duration
	dto.PatientName = session.Patient.Name
	return dto, nil
}
