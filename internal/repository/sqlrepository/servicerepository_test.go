package sqlrepository_test

import (
	"testing"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
	"github.com/AnatoliyBr/dwh-service/internal/repository/sqlrepository"
	"github.com/stretchr/testify/assert"
)

func TestServiceRepository_Create(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("services")

	s := entity.TestService(t)
	sr := sqlrepository.NewServiceRepository(db)

	assert.NoError(t, sr.Create(s))
}

func TestServiceRepository_FindByID(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("services")

	s1 := entity.TestService(t)
	sr := sqlrepository.NewServiceRepository(db)

	sr.Create(s1)

	_, err := sr.FindByID(s1.ServiceID + 1)
	assert.EqualError(t, err, repository.ErrRecordNotFound.Error())

	s2, err := sr.FindByID(s1.ServiceID)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
}
