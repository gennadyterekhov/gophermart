package server

import (
	"net/http/httptest"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/controllers"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/router"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"
	"github.com/go-chi/chi/v5"
)

type HasServer interface {
	SetServer(srv *httptest.Server)
	GetServer() *httptest.Server
}

type BaseSuiteWithServerInterface interface {
	base.BaseSuiteInterface
	HasServer
}

type BaseSuiteWithServer struct {
	base.BaseSuite
	tests.TestHTTPServer
}

func InitBaseSuiteWithServer[T BaseSuiteWithServerInterface](srv T) {
	repo := repositories.NewRepositoryMock()
	srv.SetRepository(repo)
	controllersStruct := controllers.NewControllers(config.NewConfig(), repo)
	srv.SetServer(httptest.NewServer(
		router.NewRouter(controllersStruct).Router,
	))
}

func InitBaseSuiteWithServerUsingCustomRouter[T BaseSuiteWithServerInterface](srv T, rtr *chi.Mux) {
	repo := repositories.NewRepositoryMock()
	srv.SetRepository(repo)
	srv.SetServer(httptest.NewServer(
		rtr,
	))
}

func (s *BaseSuiteWithServer) SetRepository(repo *repositories.RepositoryMock) {
	s.Repository = repo
}

func (s *BaseSuiteWithServer) GetRepository() *repositories.RepositoryMock {
	return s.Repository
}

func (s *BaseSuiteWithServer) SetServer(srv *httptest.Server) {
	s.Server = srv
}

func (s *BaseSuiteWithServer) GetServer() *httptest.Server {
	return s.Server
}
