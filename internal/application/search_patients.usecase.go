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
	Term string `json:"term"`
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

func (u *SearchPatientsUseCase) Execute(ctx context.Context, input SearchPatientsInputDTO) ([]SearchPatientsOutputDTO, error) {
	dto := []SearchPatientsOutputDTO{}
	patients, err := u.SessionRepository.SearchPatientsByName(ctx, input.Term, 10, 0)
	if err != nil {
		return []SearchPatientsOutputDTO{}, err
	}
	for _, patient := range patients {
		var phones []PhoneOutputDTO
		for _, phone := range patient.Phones {
			phones = append(phones, PhoneOutputDTO{
				Value:  phone.Value,
				IsChat: phone.IsChat,
				IsMain: phone.IsMain,
			})
		}
		dto = append(dto, SearchPatientsOutputDTO{
			ID:       patient.ID,
			Name:     patient.Name,
			Document: patient.Document,
			Email:    patient.Email,
			Phones:   phones,
		})
	}
	return dto, nil
}
