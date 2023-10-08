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

func TestAPIServer_HandleMetricCreate(t *testing.T) {
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
				"slug":        "READING_TIME_NOTE_1",
				"metric_type": "DURATION",
				"details": `The "Read Time" metric allows you to estimate the approximate amount \
				of time it will take the user to read the page from beginning to end, \
				including the content of all snippets, variables, headers and footers, if any.`,
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
				"slug":        "READING_?#@*&%!",
				"metric_type": "-____-",
				"details":     "0_0",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/metrics", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestAPIServer_HandleMetricFindByID(t *testing.T) {
	sr := testrepository.NewServiceRepository()
	mr := testrepository.NewMetricRepository()
	er := testrepository.NewEventRepository()
	uc := usecase.NewAppUseCase(sr, mr, er)
	s, _ := NewAPIServer(NewConfig(), uc)

	metric := entity.TestMetric(t)
	mr.Create(metric)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "valid",
			payload:      map[string]int{"metric_id": 1},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid symbols",
			payload:      map[string]int{"metric_id": 2},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodGet, "/metrics", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
