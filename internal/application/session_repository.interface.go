package application

import (
	"context"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
)

type ListRepositoryInput struct {
	OwnerId  string
	PageSize int
	Page     int
}

type DeleteRepositoryInput struct {
	OwnerId string
	Id      string
}

type FindPatientByNameRepositoryInput struct {
	OwnerId string
	Name    string
}

type SearchPatientsByNameRepositoryInput struct {
	Term     string
	OwnerId  string
	PageSize int
	Page     int
}

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, input DeleteRepositoryInput) error
	List(ctx context.Context, input ListRepositoryInput) ([]*domain.Session, error)
	FindPatientByName(ctx context.Context, input FindPatientByNameRepositoryInput) (*domain.Patient, error)
	CreatePatient(ctx context.Context, patient *domain.Patient) error
	SearchPatientsByName(ctx context.Context, input SearchPatientsByNameRepositoryInput) ([]*domain.Patient, error)
}
