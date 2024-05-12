package main

import (
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers"
)

func main() {
	config := getConfig()
	fmt.Printf("Server started on %v\n", config.Addr)
	err := http.ListenAndServe(config.Addr, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
