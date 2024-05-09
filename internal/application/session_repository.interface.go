package application

import domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"

type SessionRepositoryInterface interface {
	Save(session *domain.Session) error
	Delete(id string) error
	List(pageSize int, page int) ([]domain.Session, error)
}
