package application

import (
	"context"
	"testing"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchPatientsUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewSearchPatientsUseCase(mockRepo)

	patients := []*domain.Patient{
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

	input := SearchPatientsInputDTO{
		Term:     "Ali",
		PageSize: 0,
		Page:     0,
	}

	mockRepo.On("SearchPatientsByName", mock.Anything).Return(patients, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertNumberOfCalls(t, "SearchPatientsByName", 1)

	mockRepo.AssertExpectations(t)

}
