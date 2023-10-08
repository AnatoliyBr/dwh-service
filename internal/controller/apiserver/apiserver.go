package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/AnatoliyBr/dwh-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ctxKey uint8

const (
	ctxKeyRequestID ctxKey = iota
)

type apiServer struct {
	httpServer      *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	config          *Config
	logger          *logrus.Logger
	uc              usecase.UseCase
}

func NewAPIServer(config *Config, uc usecase.UseCase) (*apiServer, error) {
	s := &apiServer{
		httpServer: &http.Server{
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			Addr:         config.BindAddr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: config.ShutdownTimeout,
		config:          config,
		logger:          logrus.New(),
		uc:              uc,
	}

	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *apiServer) configureRouter() {
	r := mux.NewRouter()

	// middleware
	r.Use(s.setRequestID)
	r.Use(s.logRequest)
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	// public
	r.HandleFunc("/services", s.handleServiceCreate()).Methods(http.MethodPost)
	r.HandleFunc("/services", s.handleServiceFindByID()).Methods(http.MethodGet)

	r.HandleFunc("/metrics", s.handleMetricCreate()).Methods(http.MethodPost)
	r.HandleFunc("/metrics", s.handleMetricFindByID()).Methods(http.MethodGet)

	r.HandleFunc("/events", s.handleEventCreate()).Methods(http.MethodPost)
	r.HandleFunc("/events", s.handleGetMetricValuesForTimePeriod()).Methods(http.MethodGet)

	s.httpServer.Handler = r
}

func (s *apiServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *apiServer) StartAPIServer() {
	s.logger.Info("starting api server")

	go func() {
		s.notify <- s.httpServer.ListenAndServe()
		close(s.notify)
	}()
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpServer.Handler.ServeHTTP(w, r)
}

func (s *apiServer) Notify() <-chan error {
	return s.notify
}

func (s *apiServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

func (s *apiServer) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *apiServer) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *apiServer) handleServiceCreate() http.HandlerFunc {
	type request struct {
		Slug    string `json:"slug"`
		Details string `json:"details"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		service := &entity.Service{
			Slug:    req.Slug,
			Details: req.Details,
		}

		if err := s.uc.ServiceCreate(service); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, service)
	}
}

func (s *apiServer) handleServiceFindByID() http.HandlerFunc {
	type request struct {
		ServiceID int `json:"service_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		service, err := s.uc.ServiceFindByID(req.ServiceID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, service)
	}
}

func (s *apiServer) handleMetricCreate() http.HandlerFunc {
	type request struct {
		Slug       string `json:"slug"`
		MetricType string `json:"metric_type"`
		Details    string `json:"details"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		metric := &entity.Metric{
			Slug:       req.Slug,
			MetricType: req.MetricType,
			Details:    req.Details,
		}

		if err := s.uc.MetricCreate(metric); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, metric)
	}
}

func (s *apiServer) handleMetricFindByID() http.HandlerFunc {
	type request struct {
		MetricID int `json:"metric_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		metric, err := s.uc.MetricFindByID(req.MetricID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, metric)
	}
}

func (s *apiServer) handleEventCreate() http.HandlerFunc {
	type request struct {
		ServiceID int                   `json:"service_id"`
		Metrics   []*entity.LightMetric `json:"metrics"`
	}

	type response struct {
		Event   *entity.Event         `json:"event"`
		Metrics []*entity.LightMetric `json:"metrics"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		e := &entity.Event{
			TimeStamp: entity.CustomTime{Time: time.Now()},
			ServiceID: req.ServiceID,
		}

		if err := s.uc.EventCreate(e); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := s.uc.AddMetricsToEvent(e.EventID, req.Metrics); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resp := &response{
			Event:   e,
			Metrics: req.Metrics,
		}

		s.respond(w, r, http.StatusCreated, resp)
	}
}

func (s *apiServer) handleGetMetricValuesForTimePeriod() http.HandlerFunc {
	type request struct {
		ServiceID int                   `json:"service_id"`
		Period    [2]*entity.CustomTime `json:"period"`
		MetricID  int                   `json:"metric_id"`
	}

	type response struct {
		Request *request      `json:"request"`
		Report  []*entity.Row `json:"report"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if !req.Period[0].Time.Before(req.Period[1].Time) {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid period"))
			return
		}

		_, err := s.uc.ServiceFindByID(req.ServiceID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		metric, err := s.uc.MetricFindByID(req.MetricID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		report, err := s.uc.GetMetricValuesForTimePeriod(req.ServiceID, req.Period, metric)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resp := &response{
			Request: req,
			Report:  report.([]*entity.Row),
		}

		s.respond(w, r, http.StatusOK, resp)
	}
}

func (s *apiServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *apiServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(data)
	}
}
