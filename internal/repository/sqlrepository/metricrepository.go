package sqlrepository

import (
	"database/sql"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

type MetricRepository struct {
	db *sql.DB
}

func NewMetricRepository(db *sql.DB) *MetricRepository {
	return &MetricRepository{
		db: db,
	}
}

func (r *MetricRepository) Create(m *entity.Metric) error {
	if err := m.Validate(); err != nil {
		return err
	}

	return r.db.QueryRow(
		"INSERT INTO metrics (slug, metric_type, details) VALUES ($1, $2, $3) RETURNING metric_id",
		m.Slug,
		m.MetricType,
		m.Details,
	).Scan(&m.MetricID)
}

func (r *MetricRepository) FindByID(metricID int) (*entity.Metric, error) {
	m := &entity.Metric{}
	if err := r.db.QueryRow(
		"SELECT metric_id, slug, metric_type, details FROM metrics WHERE metric_id = $1",
		metricID,
	).Scan(
		&m.MetricID,
		&m.Slug,
		&m.MetricType,
		&m.Details,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}
	return m, nil
}
