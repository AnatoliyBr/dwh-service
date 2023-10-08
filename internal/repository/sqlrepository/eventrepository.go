package sqlrepository

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

const (
	defaultLayout = time.RFC3339
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) Create(e *entity.Event) error {
	return r.db.QueryRow(
		"INSERT INTO events (time_stamp, service_id) VALUES ($1, $2) RETURNING event_id",
		e.TimeStamp.Time,
		e.ServiceID,
	).Scan(&e.EventID)
}

func (r *EventRepository) AddMetricsToEvent(eventID int, metrics []*entity.LightMetric) error {
	stmt, err := r.db.Prepare(
		"INSERT INTO events_with_metrics (event_id, metric_id, metric_value) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	for _, m := range metrics {
		_, err := stmt.Exec(eventID, m.MetricID, m.MetricValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *EventRepository) GetMetricValuesForTimePeriod(serviceID int, p [2]*entity.CustomTime, m *entity.Metric) (interface{}, error) {
	values := make([]*entity.Row, 0)

	rows, err := r.db.Query(
		`SELECT e.time_stamp, ewm.metric_value FROM events e, events_with_metrics ewm WHERE (ewm.event_id IN (SELECT event_id FROM events WHERE service_id = $1 AND (time_stamp >= $2 AND time_stamp <= $3))) AND ewm.metric_id = $4 AND ewm.event_id = e.event_id`,
		serviceID,
		p[0].Time,
		p[1].Time,
		m.MetricID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var v string

		err := rows.Scan(&t, &v)
		if err != nil {
			return nil, err
		}

		switch m.MetricType {
		case "INT":
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			values = append(values, &entity.Row{
				TimeStamp: &entity.CustomTime{
					Time: t,
				},
				Value: i,
			})
		case "FLOAT":
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, err
			}
			values = append(values, &entity.Row{
				TimeStamp: &entity.CustomTime{
					Time: t,
				},
				Value: f,
			})
		case "DURATION":
			d, err := time.ParseDuration(v)
			if err != nil {
				return nil, err
			}
			values = append(values, &entity.Row{
				TimeStamp: &entity.CustomTime{
					Time: t,
				},
				Value: d.String(),
			})
		case "TIMESTAMP_WITH_TIMEZONE":
			tmstmp, err := time.Parse(defaultLayout, v)
			if err != nil {
				return nil, err
			}
			values = append(values, &entity.Row{
				TimeStamp: &entity.CustomTime{
					Time: t,
				},
				Value: &entity.CustomTime{Time: tmstmp},
			})
		case "BOOL":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, err
			}
			values = append(values, &entity.Row{
				TimeStamp: &entity.CustomTime{
					Time: t,
				},
				Value: b,
			})
		case "STRING":
			values = append(values, &entity.Row{
				TimeStamp: &entity.CustomTime{
					Time: t,
				},
				Value: v,
			})
		default:
			return nil, errors.New("unknown metric type")
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(values) > 0 {
		return values, nil
	} else {
		return nil, repository.ErrRecordNotFound
	}
}
