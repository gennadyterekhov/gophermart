package main

import (
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers"
)

func main() {
	fmt.Println("gophermart initialization")
	logger.Init()
	config.Init()
	fmt.Printf("Server started on %v\n", config.ServerConfig.Addr)
	err := http.ListenAndServe(config.ServerConfig.Addr, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
