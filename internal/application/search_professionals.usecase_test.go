package application

import (
	"context"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
	"testing"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchProfessionalsUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewSearchProfessionalsUseCase(mockRepo)

	professionals := []domain.Professional{
		{
			ID:   uuid.New().String(),
			Name: "Aliba",
		},
		{
			ID:   uuid.New().String(),
			Name: "Alice",
		},
		{
			ID:   uuid.New().String(),
			Name: "Aliu",
		},
	}

	input := SearchProfessionalsInputDTO{
		Term:     "Ali",
		PageSize: 0,
		Page:     0,
	}

	mockRepo.On("SearchProfessionalsByName", mock.Anything).Return(&helpers.Pages[domain.Professional]{
		Content: professionals,
	}, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertNumberOfCalls(t, "SearchProfessionalsByName", 1)

	mockRepo.AssertExpectations(t)

}
