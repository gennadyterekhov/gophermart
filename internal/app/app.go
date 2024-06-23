package app

import (
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/controllers"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/router"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

type GophermartApp struct {
	ServerConfig *config.Config
	DB           *storage.DB
	Router       *router.Router
}

func NewApp() *GophermartApp {
	serverConfig := config.NewConfig()
	db := storage.NewDB(serverConfig.DBDsn)
	repo := repositories.NewRepository(db)
	controllersStruct := controllers.NewControllers(serverConfig, repo)
	routerInstance := router.NewRouter(controllersStruct)

	app := &GophermartApp{
		ServerConfig: serverConfig,
		DB:           db,
		Router:       routerInstance,
	}

	return app
}

func (a GophermartApp) StartServer() error {
	fmt.Printf("Server started on %v\n", a.ServerConfig.Addr)
	err := http.ListenAndServe(a.ServerConfig.Addr, a.Router.Router)

	return err
}
