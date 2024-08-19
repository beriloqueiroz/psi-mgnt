package application

import (
	"context"
	"testing"
	"time"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindSessionUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewFindSessionUseCase(mockRepo)
	id := uuid.New().String()
	session := &domain.Session{
		ID:       id,
		Price:    12,
		Notes:    "nota 1",
		Date:     time.Now(),
		Duration: 10,
		Patient: &domain.Patient{
			ID:   "123",
			Name: "teste fulano",
		},
		Professional: &domain.Professional{
			ID:   "123",
			Name: "teste proff",
		},
	}

	mockRepo.On("Find", mock.Anything).Return(session, nil)

	output, err := usecase.Execute(context.Background(), FindSessionInputDto{ID: id})

	assert.NoError(t, err)
	assert.NotNil(t, output)

	mockRepo.AssertNumberOfCalls(t, "Find", 1)

	mockRepo.AssertExpectations(t)
}
