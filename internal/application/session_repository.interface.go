package application

import (
	"context"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
)

type ListByProfessionalRepositoryInput struct {
	ProfessionalId string
	ListConfig     helpers.ListConfig
}

type ListRepositoryInput struct {
	ListConfig helpers.ListConfig
}

type DeleteRepositoryInput struct {
	ProfessionalId string
	SessionId      string
}

type FindPatientRepositoryInput struct {
	PatientId string
}

type FindProfessionalRepositoryInput struct {
	ProfessionalId string
}

type SearchPatientsByNameRepositoryInput struct {
	ListConfig helpers.ListConfig
}

type SearchProfessionalByNameRepositoryInput struct {
	ListConfig helpers.ListConfig
}

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, input DeleteRepositoryInput) error
	ListByProfessional(ctx context.Context, input ListByProfessionalRepositoryInput) (*helpers.Pages[domain.Session], error)
	List(ctx context.Context, input ListRepositoryInput) (*helpers.Pages[domain.Session], error)
	FindPatient(ctx context.Context, input FindPatientRepositoryInput) (*domain.Patient, error)
	CreatePatient(ctx context.Context, patient *domain.Patient) error
	SearchPatientsByName(ctx context.Context, input SearchPatientsByNameRepositoryInput) (*helpers.Pages[domain.Patient], error)
	FindProfessional(ctx context.Context, input FindProfessionalRepositoryInput) (*domain.Professional, error)
	SearchProfessionalsByName(ctx context.Context, input SearchProfessionalByNameRepositoryInput) (*helpers.Pages[domain.Professional], error)
	CreateProfessional(ctx context.Context, professional *domain.Professional) error
}
