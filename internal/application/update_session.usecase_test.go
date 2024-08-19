package application

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSessionUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewUpdateSessionUseCase(mockRepo)

	input := UpdateSessionInputDTO{
		ID:    "123",
		Notes: "Test notes",
	}

	existingPatient := &domain.Patient{
		ID:   "123",
		Name: "John Doe",
	}

	existingProfessional := &domain.Professional{
		ID:   "123",
		Name: "John Doe Prof",
	}

	existingSession := &domain.Session{
		ID:           "123",
		Notes:        "Test notes",
		Date:         time.Now(),
		Duration:     time.Minute,
		Patient:      existingPatient,
		Professional: existingProfessional,
		Price:        50,
	}
	mockRepo.On("Find", mock.Anything).Return(existingSession, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, existingSession.Price, output.Price)
	assert.Equal(t, input.Notes, output.Notes)
	assert.Equal(t, existingSession.Date, output.Date)
	assert.Equal(t, existingSession.Duration, output.Duration)
	mockRepo.AssertNumberOfCalls(t, "Update", 1)
	mockRepo.AssertNumberOfCalls(t, "Find", 1)

	mockRepo.AssertExpectations(t)
}

func TestUpdateSessionUseCase_Execute_FindSessionError(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewUpdateSessionUseCase(mockRepo)

	input := UpdateSessionInputDTO{
		Notes: "Test notes",
		ID:    "123",
	}

	inputRepo := FindSessionRepositoryInput{
		ID: "123",
	}

	mockRepo.On("Find", inputRepo).Return(&domain.Session{}, errors.New("something went wrong"))

	_, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)

	mockRepo.AssertNumberOfCalls(t, "Update", 0)
	mockRepo.AssertNumberOfCalls(t, "Find", 1)
	mockRepo.AssertExpectations(t)
}
