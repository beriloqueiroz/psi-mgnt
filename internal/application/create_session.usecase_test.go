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

func TestCreateSessionUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateSessionUseCase(mockRepo)

	input := CreateSessionInputDTO{
		Price:          100,
		Notes:          "Test notes",
		Date:           time.Now(),
		Duration:       30 * time.Minute,
		PatientName:    "Patient Name",
		PatientId:      "123",
		ProfessionalId: "123",
	}

	existingPatient := &domain.Patient{
		ID:   "123",
		Name: "John Doe",
	}

	existingProfessional := &domain.Professional{
		ID:   "123",
		Name: "John Doe Prof",
	}
	mockRepo.On("FindPatient", mock.Anything).Return(existingPatient, nil)
	mockRepo.On("FindProfessional", mock.Anything).Return(existingProfessional, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.Price, output.Price)
	assert.Equal(t, input.Notes, output.Notes)
	assert.Equal(t, input.Date, output.Date)
	assert.Equal(t, input.Duration, output.Duration)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "FindPatient", 1)

	mockRepo.AssertExpectations(t)
}

func TestCreateSessionUseCase_Execute_NewPatient(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateSessionUseCase(mockRepo)

	input := CreateSessionInputDTO{
		Price:          100,
		Notes:          "Test notes",
		Date:           time.Now(),
		Duration:       30 * time.Minute,
		PatientName:    "John Doe",
		PatientId:      "123",
		ProfessionalId: "123",
	}

	existingProfessional := &domain.Professional{
		ID:   "123",
		Name: "John Doe Prof",
	}

	mockRepo.On("FindPatient", mock.Anything).Return(nil, nil)
	mockRepo.On("CreatePatient", mock.Anything).Return(nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	mockRepo.On("FindProfessional", mock.Anything).Return(existingProfessional, nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.Price, output.Price)
	assert.Equal(t, input.Notes, output.Notes)
	assert.Equal(t, input.Date, output.Date)
	assert.Equal(t, input.Duration, output.Duration)
	mockRepo.AssertNumberOfCalls(t, "CreatePatient", 1)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "FindPatient", 1)

	mockRepo.AssertExpectations(t)
}

func TestCreateSessionUseCase_Execute_FindPatientError(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateSessionUseCase(mockRepo)

	input := CreateSessionInputDTO{
		Price:          100,
		Notes:          "Test notes",
		Date:           time.Now(),
		Duration:       30 * time.Minute,
		PatientName:    "John Doe",
		PatientId:      "123",
		ProfessionalId: "123",
	}

	inputRepo := FindPatientRepositoryInput{
		PatientId: "123",
	}

	mockRepo.On("FindPatient", inputRepo).Return(&domain.Patient{}, errors.New("something went wrong"))

	_, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)

	mockRepo.AssertNumberOfCalls(t, "CreatePatient", 0)
	mockRepo.AssertNumberOfCalls(t, "Create", 0)
	mockRepo.AssertNumberOfCalls(t, "FindPatient", 1)
	mockRepo.AssertExpectations(t)
}
