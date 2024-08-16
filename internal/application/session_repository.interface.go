package application

import (
	"context"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
)

type ListByProfessionalRepositoryInput struct {
	ProfessionalId  string
	PageSize        int
	Page            int
	PatientNameTerm string `json:"patient_name"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
}

type ListRepositoryInput struct {
	PageSize        int
	Page            int
	PatientNameTerm string `json:"patient_name"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
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
	Term     string
	PageSize int
	Page     int
}

type SearchProfessionalByNameRepositoryInput struct {
	Term     string
	PageSize int
	Page     int
}

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, input DeleteRepositoryInput) error
	ListByProfessional(ctx context.Context, input ListByProfessionalRepositoryInput) ([]*domain.Session, error)
	List(ctx context.Context, input ListRepositoryInput) ([]*domain.Session, error)
	FindPatient(ctx context.Context, input FindPatientRepositoryInput) (*domain.Patient, error)
	CreatePatient(ctx context.Context, patient *domain.Patient) error
	SearchPatientsByName(ctx context.Context, input SearchPatientsByNameRepositoryInput) ([]*domain.Patient, error)
	FindProfessional(ctx context.Context, input FindProfessionalRepositoryInput) (*domain.Professional, error)
	SearchProfessionalsByName(ctx context.Context, input SearchProfessionalByNameRepositoryInput) ([]*domain.Professional, error)
	CreateProfessional(ctx context.Context, professional *domain.Professional) error
}
