package main

import (
	"fmt"
	"net/http"
)

func main() {
	var config *Config = getConfig()
	fmt.Printf("Server started on %v\n", config.Addr)
	err = http.ListenAndServe(config.Addr, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
