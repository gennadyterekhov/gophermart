package main

import (
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers"
)

func main() {
	fmt.Printf("Server started on %v\n", config.ServerConfig.Addr)
	err := http.ListenAndServe(config.ServerConfig.Addr, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
