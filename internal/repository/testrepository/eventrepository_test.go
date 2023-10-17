package testrepository_test

import (
	"testing"
	"time"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
	"github.com/AnatoliyBr/dwh-service/internal/repository/testrepository"
	"github.com/stretchr/testify/assert"
)

const (
	defaultLayout = time.RFC3339
)

func TestEventRepository_Create(t *testing.T) {
	e := entity.TestEvent(t)
	er := testrepository.NewEventRepository()
	assert.NoError(t, er.Create(e))
}

func TestEventRepository_AddMetricsToEvent(t *testing.T) {
	testCases := []struct {
		name        string
		metricValue interface{}
	}{
		{
			name:        "int",
			metricValue: 10,
		},
		{
			name:        "float",
			metricValue: 56.7,
		},
		{
			name:        "duration",
			metricValue: time.Duration(10 * time.Second).String(),
		},
		{
			name:        "timestamp with timezone",
			metricValue: entity.CustomTime{Time: time.Now()}.Time.Format(defaultLayout),
		},
		{
			name:        "bool",
			metricValue: true,
		},
		{
			name:        "string",
			metricValue: "starting api server",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := entity.TestMetric(t)
			e := entity.TestEvent(t)

			er := testrepository.NewEventRepository()
			err := er.AddMetricsToEvent(e.EventID, []*entity.AddMetric{
				{
					MetricID:    m.MetricID,
					MetricValue: tc.metricValue,
				}})
			assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

			er.Create(e)
			err = er.AddMetricsToEvent(e.EventID, []*entity.AddMetric{
				{
					MetricID:    m.MetricID,
					MetricValue: tc.metricValue,
				}})

			assert.NoError(t, err)
		})
	}
}

func TestEventRepository_GetMetricValuesForTimePeriod(t *testing.T) {
	testCases := []struct {
		name        string
		m           func() *entity.Metric
		metricValue interface{}
		p           [2]*entity.CustomTime
	}{
		{
			name: "int",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "INT"
				return m
			},
			metricValue: 10,
			p: [2]*entity.CustomTime{
				{Time: time.Now().AddDate(0, 0, -1)},
				{Time: time.Now().AddDate(0, 0, +1)},
			},
		},
		{
			name: "float",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "FLOAT"
				return m
			},
			metricValue: 56.7,
			p: [2]*entity.CustomTime{
				{Time: time.Now().AddDate(0, 0, -1)},
				{Time: time.Now().AddDate(0, 0, +1)},
			},
		},
		{
			name: "duration",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "DURATION"
				return m
			},
			metricValue: time.Duration(10 * time.Second),
			p: [2]*entity.CustomTime{
				{Time: time.Now().AddDate(0, 0, -1)},
				{Time: time.Now().AddDate(0, 0, +1)},
			},
		},
		{
			name: "timestamp with timezone",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "TIMESTAMP_WITH_TIMEZONE"
				return m
			},
			metricValue: entity.CustomTime{Time: time.Now()}.Time,
			p: [2]*entity.CustomTime{
				{Time: time.Now().AddDate(0, 0, -1)},
				{Time: time.Now().AddDate(0, 0, +1)},
			},
		},
		{
			name: "bool",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "BOOL"
				return m
			},
			metricValue: true,
			p: [2]*entity.CustomTime{
				{Time: time.Now().AddDate(0, 0, -1)},
				{Time: time.Now().AddDate(0, 0, +1)},
			},
		},
		{
			name: "string",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "STRING"
				return m
			},
			metricValue: "starting api server",
			p: [2]*entity.CustomTime{
				{Time: time.Now().AddDate(0, 0, -1)},
				{Time: time.Now().AddDate(0, 0, +1)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := entity.TestService(t)
			e := entity.TestEvent(t)
			m := tc.m()

			mr := testrepository.NewMetricRepository()
			mr.Create(m)

			sr := testrepository.NewServiceRepository()
			sr.Create(s)
			e.ServiceID = s.ServiceID

			er := testrepository.NewEventRepository()
			er.Create(e)

			_, err := er.GetMetricValuesForTimePeriod(e.ServiceID, tc.p, m)
			assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

			er.AddMetricsToEvent(e.EventID, []*entity.AddMetric{
				{
					MetricID:    m.MetricID,
					MetricValue: tc.metricValue,
				}})

			_, err = er.GetMetricValuesForTimePeriod(e.ServiceID, tc.p, m)
			assert.NoError(t, err)
		})
	}
}
