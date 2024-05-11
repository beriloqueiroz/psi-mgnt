package application

import (
	"context"
)

type SearchPatientsUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewSearchPatientsUseCase(
	sessionRepository SessionRepositoryInterface,
) *SearchPatientsUseCase {
	return &SearchPatientsUseCase{
		SessionRepository: sessionRepository,
	}
}

type SearchPatientsInputDTO struct {
	Term     string `json:"term"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
	OwnerId  string `json:"owner_id"`
}

type SearchPatientsOutputDTO struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Document string           `json:"document"`
	Email    string           `json:"email"`
	Phones   []PhoneOutputDTO `json:"phones"`
}

type PhoneOutputDTO struct {
	Value  string `json:"value"`
	IsChat bool   `json:"is_chat"`
	IsMain bool   `json:"is_main"`
}

func (u *SearchPatientsUseCase) Execute(ctx context.Context, input SearchPatientsInputDTO) ([]*SearchPatientsOutputDTO, error) {
	pageSizeParsed := 50
	pageParsed := 1
	if input.PageSize >= 0 {
		pageSizeParsed = input.PageSize
	}
	if input.Page >= 1 {
		pageParsed = input.Page
	}
	repoInput := SearchPatientsByNameRepositoryInput{
		Term:     input.Term,
		OwnerId:  input.OwnerId,
		PageSize: pageSizeParsed,
		Page:     pageParsed,
	}
	patients, err := u.SessionRepository.SearchPatientsByName(ctx, repoInput)
	if err != nil {
		return []*SearchPatientsOutputDTO{}, err
	}
	dto := []*SearchPatientsOutputDTO{}
	for _, patient := range patients {
		var phones []PhoneOutputDTO
		for _, phone := range patient.Phones {
			phones = append(phones, PhoneOutputDTO{
				Value:  phone.Value,
				IsChat: phone.IsChat,
				IsMain: phone.IsMain,
			})
		}
		dto = append(dto, &SearchPatientsOutputDTO{
			ID:       patient.ID,
			Name:     patient.Name,
			Document: patient.Document,
			Email:    patient.Email,
			Phones:   phones,
		})
	}
	return dto, nil
}
