package application

import (
	"context"
	"time"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateSessionUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewCreateSessionUseCase(
	sessionRepository SessionRepositoryInterface,
) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		SessionRepository: sessionRepository,
	}
}

type CreateSessionInputDTO struct {
	Price       float64       `json:"price"`
	Notes       string        `json:"notes"`
	Date        time.Time     `json:"date"`
	PaymentDate time.Time     `json:"payment_date"`
	Duration    time.Duration `json:"duration"`
	PatientName string        `json:"patient_name"`
}

type CreateSessionOutputDTO struct {
	ID          string        `json:"id"`
	Price       float64       `json:"price"`
	Notes       string        `json:"notes"`
	Date        time.Time     `json:"date"`
	PaymentDate time.Time     `json:"payment_date"`
	Duration    time.Duration `json:"duration"`
	PatientName string        `json:"patient_name"`
}

func (u *CreateSessionUseCase) Execute(ctx context.Context, input CreateSessionInputDTO) (CreateSessionOutputDTO, error) {
	dto := CreateSessionOutputDTO{}
	patient, err := u.SessionRepository.FindPatientByName(ctx, input.PatientName)
	if err != nil {
		return dto, err
	}
	if patient == nil {
		patient, err = domain.NewPatient(uuid.New().String(), input.PatientName, "", "", []domain.Phone{})
		if err != nil {
			return dto, err
		}
		err := u.SessionRepository.CreatePatient(ctx, patient)
		if err != nil {
			return dto, err
		}
	}
	session, err := domain.NewSession(uuid.New().String(), input.Price, input.Notes, input.Date, input.PaymentDate, input.Duration, patient)
	if err != nil {
		return dto, err
	}
	u.SessionRepository.Create(ctx, session)
	dto.ID = session.ID
	dto.Price = session.Price
	dto.Notes = session.Notes
	dto.Date = session.Date
	dto.PaymentDate = session.PaymentDate
	dto.Duration = session.Duration
	dto.PatientName = session.Patient.Name
	return dto, nil
}
