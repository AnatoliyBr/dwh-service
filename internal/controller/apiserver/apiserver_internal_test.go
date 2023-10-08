package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/repository/testrepository"
	"github.com/AnatoliyBr/dwh-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAPIServer_SetRequestID(t *testing.T) {
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()
	uc := usecase.NewAppUseCase(sr, mr, er)
	s, _ := NewAPIServer(NewConfig(), uc)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	s.setRequestID(handler).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
}

func TestAPIServer_HandleServiceCreate(t *testing.T) {
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()
	uc := usecase.NewAppUseCase(sr, mr, er)
	s, _ := NewAPIServer(NewConfig(), uc)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"slug": "NOTE_BOOK",
				"details": `NoteBook is the best word processing app for all your works, \
				from taking down quick notes to writing your books, \
				eBooks and organizing your documents. This app is available for iOS and Mac devices.`,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid symbols",
			payload: map[string]string{
				"slug":    "NOTE_?#@*&%!",
				"details": "0_0",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/services", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestAPIServer_HandleServiceFindByID(t *testing.T) {
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()
	uc := usecase.NewAppUseCase(sr, mr, er)
	s, _ := NewAPIServer(NewConfig(), uc)

	service := entity.TestService(t)
	sr.Create(service)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "valid",
			payload:      map[string]int{"service_id": 1},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid symbols",
			payload:      map[string]int{"service_id": 2},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodGet, "/services", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
