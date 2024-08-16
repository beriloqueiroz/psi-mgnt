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
			ID:       uuid.New().String(),
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
		},
		{
			ID:       uuid.New().String(),
			Price:    12.5,
			Notes:    "nota 2",
			Date:     time.Now(),
			Duration: 50,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste sicrano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
		{
			ID:       uuid.New().String(),
			Price:    120,
			Notes:    "nota 3",
			Date:     time.Now(),
			Duration: 100,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste fulano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
	}

	input := ListSessionsInputDto{
		PageSize: 10,
		Page:     1,
	}

	mockRepo.On("List", mock.Anything).Return(sessions, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertNumberOfCalls(t, "List", 1)
	mockRepo.AssertNumberOfCalls(t, "ListByProfessional", 0)

	mockRepo.AssertExpectations(t)
}

func TestListSessionsUseCase_WhenProfessional_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewListSessionsUseCase(mockRepo)

	sessions := []*domain.Session{
		{
			ID:       uuid.New().String(),
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
		},
		{
			ID:       uuid.New().String(),
			Price:    12.5,
			Notes:    "nota 2",
			Date:     time.Now(),
			Duration: 50,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste sicrano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
		{
			ID:       uuid.New().String(),
			Price:    120,
			Notes:    "nota 3",
			Date:     time.Now(),
			Duration: 100,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste fulano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
	}

	input := ListSessionsInputDto{
		PageSize:       10,
		Page:           1,
		ProfessionalId: "123",
	}

	mockRepo.On("ListByProfessional", mock.Anything).Return(sessions, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "teste proff", output[0].ProfessionalName)
	assert.Equal(t, "teste proff", output[1].ProfessionalName)
	mockRepo.AssertNumberOfCalls(t, "ListByProfessional", 1)

	mockRepo.AssertExpectations(t)
}

func TestListSessionsWithPaginationUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewListSessionsUseCase(mockRepo)

	sessions := []*domain.Session{
		{
			ID:       uuid.New().String(),
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
		},
		{
			ID:       uuid.New().String(),
			Price:    12.5,
			Notes:    "nota 2",
			Date:     time.Now(),
			Duration: 50,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste sicrano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
		{
			ID:       uuid.New().String(),
			Price:    120,
			Notes:    "nota 3",
			Date:     time.Now(),
			Duration: 100,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste fulano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
	}

	input := ListSessionsInputDto{
		PageSize: 2,
		Page:     1,
	}

	mockRepo.On("List", mock.Anything).Return(sessions, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, len(output), 2)
	assert.Equal(t, output[0].Notes, "nota 1")
	assert.Equal(t, output[1].Notes, "nota 2")

	input.PageSize = 1
	input.Page = 2
	mockRepo.On("List", mock.Anything).Return(sessions, nil)

	output, err = usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, len(output), 1)
	assert.Equal(t, output[0].Notes, "nota 2")

	mockRepo.AssertNumberOfCalls(t, "List", 2)
	mockRepo.AssertNumberOfCalls(t, "ListByProfessional", 0)

	mockRepo.AssertExpectations(t)
}

func TestListSessionsUseCaseByPatientTerm_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	usecase := NewListSessionsUseCase(mockRepo)

	sessions := []*domain.Session{
		{
			ID:       uuid.New().String(),
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
		},
		{
			ID:       uuid.New().String(),
			Price:    12.5,
			Notes:    "nota 2",
			Date:     time.Now(),
			Duration: 50,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste sicrano",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
		{
			ID:       uuid.New().String(),
			Price:    120,
			Notes:    "nota 3",
			Date:     time.Now(),
			Duration: 100,
			Patient: &domain.Patient{
				ID:   "123",
				Name: "teste sicrano 2",
			},
			Professional: &domain.Professional{
				ID:   "123",
				Name: "teste proff",
			},
		},
	}

	input := ListSessionsInputDto{
		PageSize:        10,
		Page:            1,
		PatientNameTerm: "sicrano",
	}

	mockRepo.On("List", mock.Anything).Return(sessions, nil)

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertNumberOfCalls(t, "List", 1)
	mockRepo.AssertNumberOfCalls(t, "ListByProfessional", 0)

	mockRepo.AssertExpectations(t)
}
