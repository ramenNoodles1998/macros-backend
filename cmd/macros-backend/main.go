package main

import (
	"fmt"
	"net/http"

	"github.com/ramenNoodles1998/macros-backend/internal/routes"
	"github.com/rs/cors"
)

func main() {

	router := routes.Router()

	port := 3030
	addr := fmt.Sprintf(":%d", port)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	fmt.Printf("Server listening on http://localhost:%s\n", addr)
    http.ListenAndServe(":3030", handler)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}