package main

import (
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/client"

	"github.com/gennadyterekhov/gophermart/internal/app"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func main() {
	fmt.Println("gophermart initialization")

	jobsChannel := make(chan *client.Job)
	appInstance := app.NewApp(jobsChannel)

	fmt.Println("gophermart initialized successfully")

	err := appInstance.StartServer()
	close(jobsChannel)
	appInstance.AccrualClient.CloseJobsChannel()
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())
		panic(err)
	}
}
