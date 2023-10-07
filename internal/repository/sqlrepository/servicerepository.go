package sqlrepository

import (
	"database/sql"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (r *ServiceRepository) Create(s *entity.Service) error {
	if err := s.Validate(); err != nil {
		return err
	}

	return r.db.QueryRow(
		"INSERT INTO services (slug, details) VALUES ($1, $2) RETURNING service_id",
		s.Slug,
		s.Details,
	).Scan(&s.ServiceID)
}

func (r *ServiceRepository) FindByID(serviceID int) (*entity.Service, error) {
	s := &entity.Service{}
	if err := r.db.QueryRow(
		"SELECT service_id, slug, details FROM services WHERE service_id = $1",
		serviceID,
	).Scan(
		&s.ServiceID,
		&s.Slug,
		&s.Details,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}
	return s, nil
}
