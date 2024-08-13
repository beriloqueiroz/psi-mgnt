package application

import (
	"context"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
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
	PageSize       int    `json:"page_size"`
	Page           int    `json:"page"`
	ProfessionalId string `json:"professional_id"`
}

type ListSessionsOutputDto struct {
	ID               string        `json:"id"`
	Price            float64       `json:"price"`
	Notes            string        `json:"notes"`
	Date             time.Time     `json:"date"`
	Duration         time.Duration `json:"duration"`
	PatientName      string        `json:"patient_name"`
	Plan             string        `json:"plan"`
	ProfessionalName string        `json:"professional_name"`
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

	var sessions []*domain.Session
	var err error

	if input.ProfessionalId != "" {
		repoInput := ListByProfessionalRepositoryInput{
			ProfessionalId: input.ProfessionalId,
			PageSize:       pageSizeParsed,
			Page:           pageParsed,
		}
		sessions, err = u.SessionRepository.ListByProfessional(ctx, repoInput)
	} else {
		repoInput := ListRepositoryInput{
			PageSize: pageSizeParsed,
			Page:     pageParsed,
		}
		sessions, err = u.SessionRepository.List(ctx, repoInput)
	}

	if err != nil {
		return []*ListSessionsOutputDto{}, err
	}
	dto := []*ListSessionsOutputDto{}
	for _, session := range sessions {
		dto = append(dto, &ListSessionsOutputDto{
			ID:               session.ID,
			Price:            session.Price,
			Notes:            session.Notes,
			Date:             session.Date,
			Duration:         session.Duration,
			PatientName:      session.Patient.Name,
			ProfessionalName: session.Professional.Name,
			Plan:             session.Plan,
		})
	}
	return dto, nil
}
