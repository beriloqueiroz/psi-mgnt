package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSessionUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewDeleteSessionUseCase(mockRepo)

	input := DeleteSessionInputDTO{
		ID: "123",
	}

	mockRepo.On("Delete", mock.Anything).Return(nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.ID, output.ID)
	mockRepo.AssertNumberOfCalls(t, "Delete", 1)

	mockRepo.AssertExpectations(t)
}
