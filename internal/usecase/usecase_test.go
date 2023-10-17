package usecase_test

import (
	"testing"
	"time"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
	"github.com/AnatoliyBr/dwh-service/internal/repository/testrepository"
	"github.com/AnatoliyBr/dwh-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

const (
	defaultLayout = time.RFC3339
)

func TestAppUseCase_ServiceCreate(t *testing.T) {
	s := entity.TestService(t)
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()

	uc := usecase.NewAppUseCase(sr, mr, er)
	assert.NoError(t, uc.ServiceCreate(s))
}

func TestAppUseCase_ServiceFindByID(t *testing.T) {
	s1 := entity.TestService(t)
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()

	uc := usecase.NewAppUseCase(sr, mr, er)
	uc.ServiceCreate(s1)

	_, err := uc.ServiceFindByID(s1.ServiceID + 1)
	assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

	s2, err := uc.ServiceFindByID(s1.ServiceID)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
}

func TestAppUseCase_MetricCreate(t *testing.T) {
	m := entity.TestMetric(t)
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()
	uc := usecase.NewAppUseCase(sr, mr, er)

	assert.NoError(t, uc.MetricCreate(m))
}

func TestAppUseCase_MetricFindByID(t *testing.T) {
	m1 := entity.TestMetric(t)
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()

	uc := usecase.NewAppUseCase(sr, mr, er)
	uc.MetricCreate(m1)

	_, err := uc.MetricFindByID(m1.MetricID + 1)
	assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

	s2, err := uc.MetricFindByID(m1.MetricID)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
}

func TestAppUseCase_EventCreate(t *testing.T) {
	e := entity.TestEvent(t)
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()
	uc := usecase.NewAppUseCase(sr, mr, er)

	assert.NoError(t, uc.EventCreate(e))
}

func TestAppUseCase_AddMetricsToEvent(t *testing.T) {
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

			sr := testrepository.NewServiceRepository()
			mr := testrepository.NewMetricRepository()
			er := testrepository.NewEventRepository()
			uc := usecase.NewAppUseCase(sr, mr, er)

			err := uc.AddMetricsToEvent(e.EventID, []*entity.AddMetric{
				{
					MetricID:    m.MetricID,
					MetricValue: tc.metricValue,
				}})
			assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

			uc.EventCreate(e)
			err = uc.AddMetricsToEvent(e.EventID, []*entity.AddMetric{
				{
					MetricID:    m.MetricID,
					MetricValue: tc.metricValue,
				}})

			assert.NoError(t, err)
		})
	}
}

func TestAppUseCase_GetMetricValuesForTimePeriod(t *testing.T) {
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

			sr := testrepository.NewServiceRepository()
			mr := testrepository.NewMetricRepository()
			er := testrepository.NewEventRepository()
			uc := usecase.NewAppUseCase(sr, mr, er)

			uc.MetricCreate(m)
			uc.ServiceCreate(s)
			e.ServiceID = s.ServiceID
			uc.EventCreate(e)

			_, err := uc.GetMetricValuesForTimePeriod(e.ServiceID, tc.p, m)
			assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

			uc.AddMetricsToEvent(e.EventID, []*entity.AddMetric{
				{
					MetricID:    m.MetricID,
					MetricValue: tc.metricValue,
				}})

			_, err = uc.GetMetricValuesForTimePeriod(e.ServiceID, tc.p, m)
			assert.NoError(t, err)
		})
	}
}
