package application

import (
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type mockSessionRepository struct {
	mock.Mock
}

func (m *mockSessionRepository) Create(session *domain.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *mockSessionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockSessionRepository) List(pageSize int, page int) ([]domain.Session, error) {
	args := m.Called(pageSize, page)
	return args.Get(0).([]domain.Session), args.Error(1)
}

func (m *mockSessionRepository) FindPatientByName(name string) (*domain.Patient, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *mockSessionRepository) CreatePatient(patient *domain.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *mockSessionRepository) SearchPatientsByName(term string, pageSize int, page int) ([]domain.Patient, error) {
	args := m.Called(term)
	return args.Get(0).([]domain.Patient), args.Error(1)
}
