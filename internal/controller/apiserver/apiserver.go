package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type apiServer struct {
	httpServer      *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	config          *Config
	logger          *logrus.Logger
}

func NewAPIServer(config *Config) (*apiServer, error) {
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
	}

	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *apiServer) configureRouter() {
	r := mux.NewRouter()

	// test
	r.HandleFunc("/hello", s.handleHello()).Methods(http.MethodGet)

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

func (s *apiServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"test": "hello"})
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
