package testrepository

import (
	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

type MetricRepository struct {
	metrics map[int]*entity.Metric
}

func NewMetricRepository() *MetricRepository {
	return &MetricRepository{
		metrics: make(map[int]*entity.Metric),
	}
}

func (r *MetricRepository) Create(m *entity.Metric) error {
	if err := m.Validate(); err != nil {
		return err
	}

	m.MetricID = len(r.metrics) + 1
	r.metrics[m.MetricID] = m

	return nil
}

func (r *MetricRepository) FindByID(metricID int) (*entity.Metric, error) {
	s, ok := r.metrics[metricID]
	if !ok {
		return nil, repository.ErrRecordNotFound
	}
	return s, nil
}
