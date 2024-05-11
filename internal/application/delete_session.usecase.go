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
	ID      string `json:"id"`
	OwnerId string `json:"ownerId"`
}

type DeleteSessionOutputDTO struct {
	ID string `json:"id"`
}

func (u *DeleteSessionUseCase) Execute(ctx context.Context, input DeleteSessionInputDTO) (DeleteSessionOutputDTO, error) {
	dto := DeleteSessionOutputDTO{}
	inputRepo := DeleteRepositoryInput{
		OwnerId: input.OwnerId,
		Id:      input.ID,
	}
	err := u.SessionRepository.Delete(ctx, inputRepo)
	if err != nil {
		return dto, err
	}
	dto.ID = input.ID
	return dto, nil
}
