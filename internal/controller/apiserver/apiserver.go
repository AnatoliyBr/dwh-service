package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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

	// middleware
	r.Use(s.setRequestID)
	r.Use(s.logRequest)
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

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
