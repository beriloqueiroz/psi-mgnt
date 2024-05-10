package application

import (
	"context"

	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
)

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, pageSize int, page int) ([]*domain.Session, error)
	FindPatientByName(ctx context.Context, name string) (*domain.Patient, error)
	CreatePatient(ctx context.Context, patient *domain.Patient) error
	SearchPatientsByName(ctx context.Context, term string, pageSize int, page int) ([]*domain.Patient, error)
}
