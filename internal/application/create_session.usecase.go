package application

import "time"

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
	ID          string        `json:"id"`
	Price       float64       `json:"price"`
	Notes       string        `json:"notes"`
	Date        time.Time     `json:"date"`
	PaymentDate time.Time     `json:"payment_date"`
	Duration    time.Duration `json:"duration"`
	PatientName float64       `json:"patient_name"`
}

type CreateSessionOutputDTO struct {
	ID          string        `json:"id"`
	Price       float64       `json:"price"`
	Notes       string        `json:"notes"`
	Date        time.Time     `json:"date"`
	PaymentDate time.Time     `json:"payment_date"`
	Duration    time.Duration `json:"duration"`
	PatientName float64       `json:"patient_name"`
}

func (u *CreateSessionUseCase) Execute(input CreateSessionInputDTO) (CreateSessionOutputDTO, error) {
	//todo see if patient already exists by name and create if not
	//todo mount session with patient return and input
	//save session
	//mount output
	dto := CreateSessionOutputDTO{}
	return dto, nil
}
