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
	// Criação do mock do repository
	mockRepo := new(mockSessionRepository)

	// Caso de uso
	useCase := NewCreateSessionUseCase(mockRepo)

	// Entrada para o caso de uso
	input := CreateSessionInputDTO{
		Price:       100,
		Notes:       "Test notes",
		Date:        time.Now(),
		PaymentDate: time.Now(),
		Duration:    30 * time.Minute,
		PatientName: "John Doe",
	}

	// Simulando busca de paciente existente
	existingPatient := &domain.Patient{
		ID:   uuid.New().String(),
		Name: "John Doe",
	}
	mockRepo.On("FindPatientByName", input.PatientName).Return(existingPatient, nil)

	// Simulando criação de sessão
	mockRepo.On("Create", mock.Anything).Return(nil)

	// Executando caso de uso
	output, err := useCase.Execute(context.Background(), input)

	// Verificando resultados
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

	// Verificando chamadas ao repositório
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
	}

	mockRepo.On("FindPatientByName", input.PatientName).Return(&domain.Patient{}, nil)

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
	}

	mockRepo.On("FindPatientByName", input.PatientName).Return(&domain.Patient{}, errors.New("something went wrong"))

	_, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)

	mockRepo.AssertNumberOfCalls(t, "CreatePatient", 0)
	mockRepo.AssertNumberOfCalls(t, "Create", 0)
	mockRepo.AssertNumberOfCalls(t, "FindPatientByName", 1)
	mockRepo.AssertExpectations(t)
}
