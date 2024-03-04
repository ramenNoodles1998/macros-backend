package main

import (
	"fmt"
	"net/http"

	"github.com/ramenNoodles1998/macros-backend/internal/routes"
)

func main() {
	router := routes.Router()

	port := 3030
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on http://localhost:%s\n", addr)
	err :=http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}