package application

import (
	"context"
	"errors"
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
	Price          float64       `json:"price"`
	Notes          string        `json:"notes"`
	Date           time.Time     `json:"date"`
	Duration       time.Duration `json:"duration"`
	PatientId      string        `json:"patient_id"`
	PatientName    string        `json:"patient_name"`
	ProfessionalId string        `json:"professional_id"`
	Plan           string        `json:"plan"`
}

type CreateSessionOutputDTO struct {
	ID               string        `json:"id"`
	Price            float64       `json:"price"`
	Notes            string        `json:"notes"`
	Date             time.Time     `json:"date"`
	Duration         time.Duration `json:"duration"`
	PatientName      string        `json:"patient_name"`
	ProfessionalName string        `json:"professional_name"`
}

func (u *CreateSessionUseCase) Execute(ctx context.Context, input CreateSessionInputDTO) (CreateSessionOutputDTO, error) {
	dto := CreateSessionOutputDTO{}
	inputRepo := FindPatientRepositoryInput{
		PatientId: input.PatientId,
	}
	var patient *domain.Patient
	var err error
	if input.PatientId != "" {
		patient, err = u.SessionRepository.FindPatient(ctx, inputRepo)
		if err != nil {
			return dto, err
		}
	}

	if patient == nil {
		patient, err = domain.NewPatient(uuid.New().String(), input.PatientName, "", "", []domain.Phone{})
		if err != nil {
			return dto, err
		}
		err = u.SessionRepository.CreatePatient(ctx, patient)
		if err != nil {
			return dto, err
		}
	}

	inputRepoProf := FindProfessionalRepositoryInput{
		ProfessionalId: input.ProfessionalId,
	}
	professional, err := u.SessionRepository.FindProfessional(ctx, inputRepoProf)
	if err != nil {
		return dto, err
	}
	if professional == nil {
		return dto, errors.New("professional does not exist")
	}

	session, err := domain.NewSession(uuid.New().String(), input.Price, input.Notes, input.Date, input.Duration, patient, input.Plan, professional)
	if err != nil {
		return dto, err
	}
	err = u.SessionRepository.Create(ctx, session)
	if err != nil {
		return CreateSessionOutputDTO{}, err
	}
	dto.ID = session.ID
	dto.Price = session.Price
	dto.Notes = session.Notes
	dto.Date = session.Date
	dto.Duration = session.Duration
	dto.PatientName = session.Patient.Name
	dto.ProfessionalName = professional.Name
	return dto, nil
}
