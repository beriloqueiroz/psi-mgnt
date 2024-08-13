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
	Duration    time.Duration `json:"duration"`
	PatientName string        `json:"patient_name"`
	Plan        string        `json:"plan"`
	OwnerId     string        `json:"owner_id"`
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
		ProfessionalId: input.OwnerId,
		PageSize:       pageSizeParsed,
		Page:           pageParsed,
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
			Duration:    session.Duration,
			PatientName: session.Patient.Name,
			OwnerId:     input.OwnerId,
			Plan:        session.Plan,
		})
	}
	return dto, nil
}
