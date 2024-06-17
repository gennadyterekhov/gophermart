package main

import (
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/app"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

func main() {
	fmt.Println("gophermart initialization")

	appInstance := app.NewApp()

	fmt.Println("gophermart initialized successfully")

	err := appInstance.StartServer()
	if err != nil {
		logger.CustomLogger.Errorln(err.Error())
		panic(err)
	}
}
