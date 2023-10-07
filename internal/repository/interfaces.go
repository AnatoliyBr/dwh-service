package repository

import "github.com/AnatoliyBr/dwh-service/internal/entity"

type ServiceRepository interface {
	Create(*entity.Service) error
	FindByID(int) (*entity.Service, error)
}

type MetricRepository interface {
	Create(*entity.Metric) error
	FindByID(int) (*entity.Metric, error)
}
type EventRepository interface {
	Create(*entity.Event) error
	AddMetricsToEvent(int, []*entity.LightMetric) error
	GetMetricValuesForTimePeriod(int, [2]*entity.CustomTime, *entity.Metric) (interface{}, error)
}
