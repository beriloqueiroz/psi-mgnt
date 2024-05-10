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

func (m *mockSessionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockSessionRepository) List(ctx context.Context, pageSize int, page int) ([]*domain.Session, error) {
	args := m.Called(pageSize, page)
	return args.Get(0).([]*domain.Session), args.Error(1)
}

func (m *mockSessionRepository) FindPatientByName(ctx context.Context, name string) (*domain.Patient, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *mockSessionRepository) CreatePatient(ctx context.Context, patient *domain.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *mockSessionRepository) SearchPatientsByName(ctx context.Context, term string, pageSize int, page int) ([]*domain.Patient, error) {
	args := m.Called(term)
	return args.Get(0).([]*domain.Patient), args.Error(1)
}
