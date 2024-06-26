package application

import (
	"context"
	"time"
)

type ListSessionsUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewListSessionsUseCase(
	sessionRepository SessionRepositoryInterface,
) *ListSessionsUseCase {
	return &ListSessionsUseCase{
		SessionRepository: sessionRepository,
	}
}

type ListSessionsInputDto struct {
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
	OwnerId  string `json:"owner_id"`
}

type ListSessionsOutputDto struct {
	ID          string        `json:"id"`
	Price       float64       `json:"price"`
	Notes       string        `json:"notes"`
	Date        time.Time     `json:"date"`
	PaymentDate time.Time     `json:"payment_date"`
	Duration    time.Duration `json:"duration"`
	PatientName string        `json:"patient_name"`
}

func (u *ListSessionsUseCase) Execute(ctx context.Context, input ListSessionsInputDto) ([]*ListSessionsOutputDto, error) {
	pageSizeParsed := 50
	pageParsed := 1
	if input.PageSize >= 0 {
		pageSizeParsed = input.PageSize
	}
	if input.Page >= 1 {
		pageParsed = input.Page
	}
	repoInput := ListRepositoryInput{
		OwnerId:  input.OwnerId,
		PageSize: pageSizeParsed,
		Page:     pageParsed,
	}
	sessions, err := u.SessionRepository.List(ctx, repoInput)
	if err != nil {
		return []*ListSessionsOutputDto{}, err
	}
	dto := []*ListSessionsOutputDto{}
	for _, session := range sessions {
		dto = append(dto, &ListSessionsOutputDto{
			ID:          session.ID,
			Price:       session.Price,
			Notes:       session.Notes,
			Date:        session.Date,
			PaymentDate: session.PaymentDate,
			Duration:    session.Duration,
			PatientName: session.Patient.Name,
		})
	}
	return dto, nil
}
