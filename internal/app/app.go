package app

import (
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/domain/services"

	"github.com/gennadyterekhov/gophermart/internal/client"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/controllers"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/router"

	"github.com/gennadyterekhov/gophermart/internal/repositories"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

type GophermartApp struct {
	ServerConfig  *config.Config
	DB            *storage.DB
	Router        *router.Router
	AccrualClient *client.AccrualClient
}

func NewApp(jobsChannel chan *client.Job) *GophermartApp {
	app := &GophermartApp{}

	serverConfig := config.NewConfig()

	db := storage.NewDB(serverConfig.DBDsn)
	repo := repositories.NewRepository(db)
	app.AccrualClient = client.New(serverConfig.AccrualURL, repo, jobsChannel)

	servs := services.New(repo, app.AccrualClient)
	controllersStruct := controllers.NewControllers(servs)
	routerInstance := router.NewRouter(controllersStruct)

	app.ServerConfig = serverConfig
	app.DB = db
	app.Router = routerInstance

	return app
}

func (a GophermartApp) StartServer() error {
	fmt.Printf("Server started on %v\n", a.ServerConfig.Addr)
	err := http.ListenAndServe(a.ServerConfig.Addr, a.Router.Router)

	return err
}
