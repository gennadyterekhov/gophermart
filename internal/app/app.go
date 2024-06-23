package app

import (
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers"
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

type GophermartApp struct {
	ServerConfig *config.Config
	DB           *storage.DB
	Router       *handlers.Router
}

func NewApp() *GophermartApp {
	serverConfig := config.NewConfig()
	db := storage.NewDB(serverConfig.DBDsn)
	router := handlers.NewRouter(serverConfig, db)

	app := &GophermartApp{
		ServerConfig: serverConfig,
		DB:           db,
		Router:       router,
	}
	router.InitializeRoutes()

	return app
}

func (a GophermartApp) StartServer() error {
	fmt.Printf("Server started on %v\n", a.ServerConfig.Addr)
	err := http.ListenAndServe(a.ServerConfig.Addr, a.Router.Router)

	return err
}
