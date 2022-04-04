package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/app"
	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/config"
	"github.com/gorilla/mux"
)

type Server struct {
	log    Logger
	app    app.App
	socket string
	server *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
}

var ErrServerHasNotBeenStarted = errors.New("сервер не был запущен ранее")

func NewServer(logger Logger, app *app.App, config config.ServerConf) *Server {
	return &Server{
		log:    logger,
		app:    *app,
		socket: net.JoinHostPort(config.Host, config.Port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/hello", s.HelloHandler).Methods("GET")
	router.Use(s.loggingMiddleware)

	s.server = &http.Server{
		Addr:         s.socket,
		Handler:      router,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return ErrServerHasNotBeenStarted
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello"))
}
