package application

import (
	"context"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateProfessionalUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateProfessionalUseCase(mockRepo)

	input := CreateProfessionalInputDTO{
		Email:    "berilo@gmail.com",
		Document: "12365478",
		Name:     "teste professional",
		Phone:    "111225525252",
	}

	mockRepo.On("SearchProfessionalsByName", mock.Anything).Return([]*domain.Professional{}, nil)
	mockRepo.On("CreateProfessional", mock.Anything).Return(nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertNumberOfCalls(t, "SearchProfessionalsByName", 1)
	mockRepo.AssertNumberOfCalls(t, "CreateProfessional", 1)

	mockRepo.AssertExpectations(t)
}

func TestCreateProfessional_WhenExistsProfessional_UseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateProfessionalUseCase(mockRepo)

	input := CreateProfessionalInputDTO{
		Email:    "berilo@gmail.com",
		Document: "12365478",
		Name:     "teste professional",
		Phone:    "111225525252",
	}

	existingProfessional := &domain.Professional{
		ID:   "123",
		Name: "John Doe Prof",
	}
	mockRepo.On("SearchProfessionalsByName", mock.Anything).Return([]*domain.Professional{existingProfessional}, nil)

	_, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)
	mockRepo.AssertNumberOfCalls(t, "SearchProfessionalsByName", 1)
	mockRepo.AssertNumberOfCalls(t, "CreateProfessional", 0)

	mockRepo.AssertExpectations(t)
}
