package application

import (
	"context"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
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
	ListConfig     helpers.ListConfig
	ProfessionalId string `json:"professional_id"`
}

type ListSessionsOutputDto helpers.Pages[*ListSessionsOutputDtoItem]

type ListSessionsOutputDtoItem struct {
	ID               string        `json:"id"`
	Price            float64       `json:"price"`
	Notes            string        `json:"notes"`
	Date             time.Time     `json:"date"`
	Duration         time.Duration `json:"duration"`
	PatientName      string        `json:"patient_name"`
	Plan             string        `json:"plan"`
	ProfessionalName string        `json:"professional_name"`
}

func (u *ListSessionsUseCase) Execute(ctx context.Context, input ListSessionsInputDto) (ListSessionsOutputDto, error) {
	pageSizeParsed := 50
	pageParsed := 1
	if input.ListConfig.PageSize >= 0 {
		pageSizeParsed = input.ListConfig.PageSize
	}
	if input.ListConfig.Page >= 1 {
		pageParsed = input.ListConfig.Page
	}

	var sessions *helpers.Pages[domain.Session]
	var err error

	listConfig := helpers.ListConfig{SortField: "", IsAscending: true, PageSize: pageSizeParsed, Page: pageParsed, AndLogic: false, ExpressionFilters: []helpers.ExpressionFilter{}}
	if input.ProfessionalId != "" {
		repoInput := ListByProfessionalRepositoryInput{
			ProfessionalId: input.ProfessionalId,
			ListConfig:     listConfig,
		}
		sessions, err = u.SessionRepository.ListByProfessional(ctx, repoInput)
	} else {
		repoInput := ListRepositoryInput{
			listConfig,
		}
		sessions, err = u.SessionRepository.List(ctx, repoInput)
	}

	if err != nil {
		return ListSessionsOutputDto{}, err
	}
	dto := ListSessionsOutputDto{}
	for _, session := range sessions.Content {
		var profesionalName string
		if session.Professional != nil {
			profesionalName = session.Professional.Name
		}
		dto.Content = append(dto.Content, &ListSessionsOutputDtoItem{
			ID:               session.ID,
			Price:            session.Price,
			Notes:            session.Notes,
			Date:             session.Date,
			Duration:         session.Duration,
			PatientName:      session.Patient.Name,
			ProfessionalName: profesionalName,
			Plan:             session.Plan,
		})
	}
	dto.Size = sessions.Size
	dto.PageSize = sessions.PageSize
	dto.Page = sessions.Page
	return dto, nil
}
