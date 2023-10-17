package testrepository

import (
	"errors"
	"time"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
)

type Pair struct {
	eventID  int
	metricID int
}

type EventRepository struct {
	events            map[int]*entity.Event
	eventsWithMetrics map[Pair]interface{}
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events:            make(map[int]*entity.Event),
		eventsWithMetrics: make(map[Pair]interface{}),
	}
}

func (r *EventRepository) Create(e *entity.Event) error {
	e.EventID = len(r.events) + 1
	r.events[e.EventID] = e

	return nil
}

func (r *EventRepository) AddMetricsToEvent(eventID int, metrics []*entity.AddMetric) error {
	if _, ok := r.events[eventID]; !ok {
		return repository.ErrRecordNotFound
	}

	for _, m := range metrics {
		r.eventsWithMetrics[Pair{eventID: eventID, metricID: m.MetricID}] = m.MetricValue
	}

	return nil
}

func (r *EventRepository) GetMetricValuesForTimePeriod(serviceID int, p [2]*entity.CustomTime, m *entity.Metric) (interface{}, error) {
	values := make([]*entity.GetMetric, 0)

	suitableEvents := make([]*entity.Event, 0)

	for _, e := range r.events {
		if e.ServiceID == serviceID && p[0].Before(e.TimeStamp.Time) && p[1].After(e.TimeStamp.Time) {
			suitableEvents = append(suitableEvents, e)
		}
	}

	for _, se := range suitableEvents {
		v, _ := r.eventsWithMetrics[Pair{eventID: se.EventID, metricID: m.MetricID}]

		switch v.(type) {
		case int:
			values = append(values, &entity.GetMetric{
				TimeStamp: se.TimeStamp,
				Value:     v.(int),
			})
		case float64:
			values = append(values, &entity.GetMetric{
				TimeStamp: se.TimeStamp,
				Value:     v.(float64),
			})
		case time.Duration:
			values = append(values, &entity.GetMetric{
				TimeStamp: se.TimeStamp,
				Value:     v.(time.Duration),
			})
		case time.Time:
			values = append(values, &entity.GetMetric{
				TimeStamp: se.TimeStamp,
				Value:     &entity.CustomTime{Time: v.(time.Time)},
			})
		case bool:
			values = append(values, &entity.GetMetric{
				TimeStamp: se.TimeStamp,
				Value:     v.(bool),
			})
		case string:
			values = append(values, &entity.GetMetric{
				TimeStamp: se.TimeStamp,
				Value:     v.(string),
			})
		case nil:
			return nil, repository.ErrRecordNotFound
		default:
			return nil, errors.New("unknown metric type")
		}
	}

	if len(values) > 0 {
		return values, nil
	} else {
		return nil, repository.ErrRecordNotFound
	}
}
