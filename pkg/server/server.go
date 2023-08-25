// Package server contains ...
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gerrors "gitlab.com/tariandev_intelops/sample-http-service/pkg/errors"

	"gitlab.com/tariandev_intelops/sample-http-service/pkg/configuration"
	"gitlab.com/tariandev_intelops/sample-http-service/pkg/controller"

	"github.com/connect2naga/logger/logging"
	"github.com/gorilla/mux"
)

/*
Author : Nagarjuna S
Date : 30-04-2022 18:06
Project : sample-http-service
File : server.go
*/

type Server struct {
	ctx    context.Context
	logger logging.Logger
	e      *controller.EndpointHandler
	conf   *configuration.ServiceConfigurations
}

func NewServer(logger logging.Logger) (*Server, error) {

	config, err := configuration.GetServiceConfigurations()
	if err != nil {
		return nil, gerrors.NewFromError(gerrors.ServiceSetup, err)
	}

	return &Server{
		logger: logger,
		e:      controller.NewEndpointHandler(logger),
		conf:   config,
		ctx:    context.Background(),
	}, nil
}

func (s *Server) Start() {

	httpServer := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", s.conf.Port),
		Handler: s.handleRequests(),
	}

	s.logger.Infof(s.ctx, "Starting server at %v", httpServer.Addr)
	go func() {
		// err := httpServer.ListenAndServeTLS("./hack/localhost.crt", "./hack/localhost.key")
		// if err != nil && errors.Is(err, http.ErrServerClosed) {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf(s.ctx, "Unexpected server close: %v", err)
			os.Exit(1)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	s.logger.Info(s.ctx, "Closing the service gracefully")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Errorf(s.ctx, "Could not close the service gracefully: %v", err)
	}
	s.logger.Infof(s.ctx, "Service closed..")
}

func (s *Server) handleRequests() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/employees", s.e.GetAllEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", s.e.GetAllEmployeeById).Methods("GET")
	router.HandleFunc("/post", s.e.PutEmployees).Methods("POST")
	router.HandleFunc("/put", s.e.PutEmployees).Methods("PUT")
	router.HandleFunc("/delete", s.e.DeleteEmp).Methods("DELETE")
	return router
}
