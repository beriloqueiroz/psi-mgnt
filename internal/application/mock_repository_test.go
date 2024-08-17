package application

import (
	"context"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"

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

func (m *mockSessionRepository) ListByProfessional(ctx context.Context, input ListByProfessionalRepositoryInput) (*helpers.Pages[domain.Session], error) {
	args := m.Called(input)
	list := args.Get(0).(*helpers.Pages[domain.Session])
	start, end := Paginate(input.ListConfig.Page, input.ListConfig.PageSize, len(list.Content))
	return &helpers.Pages[domain.Session]{Content: list.Content[start:end]}, args.Error(1)
}

func (m *mockSessionRepository) List(ctx context.Context, input ListRepositoryInput) (*helpers.Pages[domain.Session], error) {
	args := m.Called(input)
	list := args.Get(0).(*helpers.Pages[domain.Session])
	start, end := Paginate(input.ListConfig.Page, input.ListConfig.PageSize, len(list.Content))
	return &helpers.Pages[domain.Session]{Content: list.Content[start:end]}, args.Error(1)
}

func (m *mockSessionRepository) FindPatient(ctx context.Context, input FindPatientRepositoryInput) (*domain.Patient, error) {
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

func (m *mockSessionRepository) SearchPatientsByName(ctx context.Context, input SearchPatientsByNameRepositoryInput) (*helpers.Pages[domain.Patient], error) {
	args := m.Called(input)
	return args.Get(0).(*helpers.Pages[domain.Patient]), args.Error(1)
}

func (m *mockSessionRepository) FindProfessional(ctx context.Context, input FindProfessionalRepositoryInput) (*domain.Professional, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Professional), args.Error(1)
}

func (m *mockSessionRepository) SearchProfessionalsByName(ctx context.Context, input SearchProfessionalByNameRepositoryInput) (*helpers.Pages[domain.Professional], error) {
	args := m.Called(input)
	return args.Get(0).(*helpers.Pages[domain.Professional]), args.Error(1)
}

func (m *mockSessionRepository) CreateProfessional(ctx context.Context, professional *domain.Professional) error {
	args := m.Called(professional)
	return args.Error(0)
}

func Paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := (pageNum - 1) * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
