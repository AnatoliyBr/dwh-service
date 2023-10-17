package usecase

import (
	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

type AppUseCase struct {
	serviceRepository repository.ServiceRepository
	metricRepository  repository.MetricRepository
	eventRepository   repository.EventRepository
}

func NewAppUseCase(sr repository.ServiceRepository, mr repository.MetricRepository, er repository.EventRepository) *AppUseCase {
	return &AppUseCase{
		serviceRepository: sr,
		metricRepository:  mr,
		eventRepository:   er,
	}
}

func (uc *AppUseCase) ServiceCreate(s *entity.Service) error {
	return uc.serviceRepository.Create(s)
}

func (uc *AppUseCase) ServiceFindByID(serviceID int) (*entity.Service, error) {
	return uc.serviceRepository.FindByID(serviceID)
}

func (uc *AppUseCase) MetricCreate(m *entity.Metric) error {
	return uc.metricRepository.Create(m)
}

func (uc *AppUseCase) MetricFindByID(metricID int) (*entity.Metric, error) {
	return uc.metricRepository.FindByID(metricID)
}

func (uc *AppUseCase) EventCreate(e *entity.Event) error {
	return uc.eventRepository.Create(e)
}

func (uc *AppUseCase) AddMetricsToEvent(eventID int, metrics []*entity.AddMetric) error {
	return uc.eventRepository.AddMetricsToEvent(eventID, metrics)
}

func (uc *AppUseCase) GetMetricValuesForTimePeriod(serviceID int, p [2]*entity.CustomTime, m *entity.Metric) (interface{}, error) {
	return uc.eventRepository.GetMetricValuesForTimePeriod(serviceID, p, m)
}
