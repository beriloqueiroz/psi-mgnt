package application

import (
	"context"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type mockSessionRepository struct {
	mock.Mock
}

func (m *mockSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *mockSessionRepository) Delete(ctx context.Context, input DeleteRepositoryInput) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockSessionRepository) List(ctx context.Context, input ListRepositoryInput) ([]*domain.Session, error) {
	args := m.Called(input)
	return args.Get(0).([]*domain.Session), args.Error(1)
}

func (m *mockSessionRepository) FindPatientByName(ctx context.Context, input FindPatientByNameRepositoryInput) (*domain.Patient, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *mockSessionRepository) CreatePatient(ctx context.Context, patient *domain.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *mockSessionRepository) SearchPatientsByName(ctx context.Context, input SearchPatientsByNameRepositoryInput) ([]*domain.Patient, error) {
	args := m.Called(input)
	return args.Get(0).([]*domain.Patient), args.Error(1)
}
