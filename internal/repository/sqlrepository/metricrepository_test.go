package sqlrepository_test

import (
	"testing"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
	"github.com/AnatoliyBr/dwh-service/internal/repository/sqlrepository"
	"github.com/stretchr/testify/assert"
)

func TestMetricRepository_Create(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("metrics")

	m := entity.TestMetric(t)
	mr := sqlrepository.NewMetricRepository(db)

	assert.NoError(t, mr.Create(m))
}

func TestMetricRepository_FindByID(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("metrics")

	m1 := entity.TestMetric(t)
	mr := sqlrepository.NewMetricRepository(db)

	mr.Create(m1)

	_, err := mr.FindByID(m1.MetricID + 1)
	assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

	s2, err := mr.FindByID(m1.MetricID)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
}
