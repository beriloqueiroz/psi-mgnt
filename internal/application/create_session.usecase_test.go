package application

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSessionUseCase_Execute(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateSessionUseCase(mockRepo)

	input := CreateSessionInputDTO{
		Price:       100,
		Notes:       "Test notes",
		Date:        time.Now(),
		PaymentDate: time.Now(),
		Duration:    30 * time.Minute,
		PatientName: "John Doe",
		OwnerId:     "123",
	}

	existingPatient := &domain.Patient{
		ID:   uuid.New().String(),
		Name: "John Doe",
	}
	mockRepo.On("FindPatientByName", mock.Anything).Return(existingPatient, nil)

	mockRepo.On("Create", mock.Anything).Return(nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.Price, output.Price)
	assert.Equal(t, input.Notes, output.Notes)
	assert.Equal(t, input.Date, output.Date)
	assert.Equal(t, input.PaymentDate, output.PaymentDate)
	assert.Equal(t, input.Duration, output.Duration)
	assert.Equal(t, input.PatientName, output.PatientName)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "FindPatientByName", 1)

	mockRepo.AssertExpectations(t)
}

func TestCreateSessionUseCase_Execute_NewPatient(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateSessionUseCase(mockRepo)

	input := CreateSessionInputDTO{
		Price:       100,
		Notes:       "Test notes",
		Date:        time.Now(),
		PaymentDate: time.Now(),
		Duration:    30 * time.Minute,
		PatientName: "John Doe",
		OwnerId:     "123",
	}

	mockRepo.On("FindPatientByName", mock.Anything).Return(nil, nil)

	mockRepo.On("CreatePatient", mock.Anything).Return(nil)

	mockRepo.On("Create", mock.Anything).Return(nil)

	output, err := useCase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.Price, output.Price)
	assert.Equal(t, input.Notes, output.Notes)
	assert.Equal(t, input.Date, output.Date)
	assert.Equal(t, input.PaymentDate, output.PaymentDate)
	assert.Equal(t, input.Duration, output.Duration)
	assert.Equal(t, input.PatientName, output.PatientName)
	mockRepo.AssertNumberOfCalls(t, "CreatePatient", 1)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "FindPatientByName", 1)

	mockRepo.AssertExpectations(t)
}

func TestCreateSessionUseCase_Execute_FindPatientError(t *testing.T) {
	mockRepo := new(mockSessionRepository)

	useCase := NewCreateSessionUseCase(mockRepo)

	input := CreateSessionInputDTO{
		Price:       100,
		Notes:       "Test notes",
		Date:        time.Now(),
		PaymentDate: time.Now(),
		Duration:    30 * time.Minute,
		PatientName: "John Doe",
		OwnerId:     "123",
	}

	inputRepo := FindPatientByNameRepositoryInput{
		OwnerId: "123",
		Name:    input.PatientName,
	}

	mockRepo.On("FindPatientByName", inputRepo).Return(&domain.Patient{}, errors.New("something went wrong"))

	_, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)

	mockRepo.AssertNumberOfCalls(t, "CreatePatient", 0)
	mockRepo.AssertNumberOfCalls(t, "Create", 0)
	mockRepo.AssertNumberOfCalls(t, "FindPatientByName", 1)
	mockRepo.AssertExpectations(t)
}
