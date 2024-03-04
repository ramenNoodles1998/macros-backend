package routes

import (
	"fmt"
	"net/http"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/data", apiDataHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from index")
}


func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	data := "data from api"
	fmt.Fprintln(w, data)
}