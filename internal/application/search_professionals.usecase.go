package application

import (
	"context"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
)

type SearchProfessionalsUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewSearchProfessionalsUseCase(
	sessionRepository SessionRepositoryInterface,
) *SearchProfessionalsUseCase {
	return &SearchProfessionalsUseCase{
		SessionRepository: sessionRepository,
	}
}

type SearchProfessionalsInputDTO struct {
	Term     string `json:"term"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

type SearchProfessionalsOutputDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Document string `json:"document"`
	Email    string `json:"email"`
}

func (u *SearchProfessionalsUseCase) Execute(ctx context.Context, input SearchProfessionalsInputDTO) ([]*SearchProfessionalsOutputDTO, error) {
	pageSizeParsed := 50
	pageParsed := 1
	if input.PageSize >= 0 {
		pageSizeParsed = input.PageSize
	}
	if input.Page >= 1 {
		pageParsed = input.Page
	}
	listConfig := helpers.ListConfig{
		PageSize: pageSizeParsed,
		Page:     pageParsed,
		ExpressionFilters: []helpers.ExpressionFilter{
			{PropertyName: "name", Value: input.Term},
		},
	}
	repoInput := SearchProfessionalByNameRepositoryInput{
		ListConfig: listConfig,
	}
	patients, err := u.SessionRepository.SearchProfessionalsByName(ctx, repoInput)
	if err != nil {
		return []*SearchProfessionalsOutputDTO{}, err
	}
	dto := []*SearchProfessionalsOutputDTO{}
	for _, patient := range patients.Content {
		dto = append(dto, &SearchProfessionalsOutputDTO{
			ID:       patient.ID,
			Name:     patient.Name,
			Document: patient.Document,
			Email:    patient.Email,
		})
	}
	return dto, nil
}
