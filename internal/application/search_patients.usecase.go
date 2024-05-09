package application

import "context"

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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Document string `json:"document"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (u *SearchPatientsUseCase) Execute(ctx context.Context, input SearchPatientsInputDTO) (SearchPatientsOutputDTO, error) {
	dto := SearchPatientsOutputDTO{}
	return dto, nil
}
