package usecase

import "github.com/AnatoliyBr/dwh-service/internal/entity"

type UseCase interface {
	ServiceCreate(*entity.Service) error
	ServiceFindByID(int) (*entity.Service, error)

	MetricCreate(*entity.Metric) error
	MetricFindByID(int) (*entity.Metric, error)

	EventCreate(*entity.Event) error
	AddMetricsToEvent(int, []*entity.AddMetric) error
	GetMetricValuesForTimePeriod(int, [2]*entity.CustomTime, *entity.Metric) (interface{}, error)
}
