package application

import domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"

type SessionRepositoryInterface interface {
	Create(session *domain.Session) error
	Delete(id string) error
	List(pageSize int, page int) ([]domain.Session, error)
	FindPatientByName(name string) (*domain.Patient, error)
	CreatePatient(patient *domain.Patient) error
	SearchPatientsByName(term string, pageSize int, page int) ([]domain.Patient, error)
}
