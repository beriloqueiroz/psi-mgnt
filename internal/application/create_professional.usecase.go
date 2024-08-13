package application

import (
	"context"
	"errors"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateProfessionalUseCase struct {
	SessionRepository SessionRepositoryInterface
}

func NewCreateProfessionalUseCase(
	sessionRepository SessionRepositoryInterface,
) *CreateProfessionalUseCase {
	return &CreateProfessionalUseCase{
		SessionRepository: sessionRepository,
	}
}

type CreateProfessionalInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Phone    string `json:"phone"`
}

type CreateProfessionalOutputDTO struct {
	ID string `json:"id"`
}

func (u *CreateProfessionalUseCase) Execute(ctx context.Context, input CreateProfessionalInputDTO) (CreateProfessionalOutputDTO, error) {
	dto := CreateProfessionalOutputDTO{}
	professionalsFound, err := u.SessionRepository.SearchProfessionalsByName(ctx, SearchProfessionalByNameRepositoryInput{
		PageSize: 1,
		Page:     1,
		Term:     input.Name,
	})
	if err != nil {
		return dto, err
	}
	if len(professionalsFound) > 0 {
		return dto, errors.New("professional name already exists")
	}
	professional, err := domain.NewProfessional(uuid.New().String(), input.Name, input.Document, input.Email)
	if err != nil {
		return dto, err
	}
	err = u.SessionRepository.CreateProfessional(ctx, professional)
	if err != nil {
		return dto, err
	}
	dto.ID = professional.ID
	return dto, nil
}
