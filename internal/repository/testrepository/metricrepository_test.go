package testrepository_test

import (
	"testing"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
	"github.com/AnatoliyBr/dwh-service/internal/repository/testrepository"
	"github.com/stretchr/testify/assert"
)

func TestMetricRepository_Create(t *testing.T) {
	m := entity.TestMetric(t)
	mr := testrepository.NewMetricRepository()

	assert.NoError(t, mr.Create(m))
}

func TestMetricRepository_FindByID(t *testing.T) {
	m1 := entity.TestMetric(t)
	mr := testrepository.NewMetricRepository()

	mr.Create(m1)

	_, err := mr.FindByID(m1.MetricID + 1)
	assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

	s2, err := mr.FindByID(m1.MetricID)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
}
