package testrepository

import (
	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

type ServiceRepository struct {
	services map[int]*entity.Service
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		services: make(map[int]*entity.Service),
	}
}

func (r *ServiceRepository) Create(s *entity.Service) error {
	if err := s.Validate(); err != nil {
		return err
	}

	s.ServiceID = len(r.services) + 1
	r.services[s.ServiceID] = s

	return nil
}

func (r *ServiceRepository) FindByID(serviceID int) (*entity.Service, error) {
	s, ok := r.services[serviceID]
	if !ok {
		return nil, repository.ErrRecordNotFound
	}
	return s, nil
}
