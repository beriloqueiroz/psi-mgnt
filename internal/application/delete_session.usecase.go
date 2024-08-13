package application

import (
	"context"
)

type DeleteSessionUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewDeleteSessionUseCase(
	sessionRepository SessionRepositoryInterface,
) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{
		SessionRepository: sessionRepository,
	}
}

type DeleteSessionInputDTO struct {
	ID string `json:"id"`
}

type DeleteSessionOutputDTO struct {
	ID string `json:"id"`
}

func (u *DeleteSessionUseCase) Execute(ctx context.Context, input DeleteSessionInputDTO) (DeleteSessionOutputDTO, error) {
	dto := DeleteSessionOutputDTO{}
	inputRepo := DeleteRepositoryInput{
		SessionId: input.ID,
	}
	err := u.SessionRepository.Delete(ctx, inputRepo)
	if err != nil {
		return dto, err
	}
	dto.ID = input.ID
	return dto, nil
}
