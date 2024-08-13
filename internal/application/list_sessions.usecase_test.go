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

func TestListSessionsUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewListSessionsUseCase(mockRepo)

	sessions := []*domain.Session{
		{
			ID:          uuid.New().String(),
			Price:       12,
			Notes:       "nota 1",
			Date:        time.Now(),
			PaymentDate: time.Now(),
			Duration:    10,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste fulano",
			},
		},
		{
			ID:          uuid.New().String(),
			Price:       12.5,
			Notes:       "nota 2",
			Date:        time.Now(),
			PaymentDate: time.Now(),
			Duration:    50,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste sicrano",
			},
		},
		{
			ID:          uuid.New().String(),
			Price:       120,
			Notes:       "nota 3",
			Date:        time.Now(),
			PaymentDate: time.Now(),
			Duration:    100,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste fulano",
			},
		},
	}

	input := ListSessionsInputDto{
		PageSize: 10,
		Page:     1,
		OwnerId:  "123",
	}

	mockRepo.On("List", mock.Anything).Return(sessions, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertNumberOfCalls(t, "List", 1)

	mockRepo.AssertExpectations(t)

}
