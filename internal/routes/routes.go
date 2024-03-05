package routes

import (
	"fmt"
	"net/http"

	macrolog "github.com/ramenNoodles1998/macros-backend/internal/macro-log"
)
func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/add-macro-log", addMacroLog)
	mux.HandleFunc("/api/get-macro-log", getMacroLog)
	mux.HandleFunc("/api/get-macro-log-by-id", getMacroLogId)
	mux.HandleFunc("/api/delete-macro-log", deleteMacroLog)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from index")
}

func addMacroLog(w http.ResponseWriter, r *http.Request) {
	macrolog.SaveMacroLog(w, r)
}

func getMacroLog(w http.ResponseWriter, r *http.Request) {
	macrolog.GetMacroLog(w, r)
}


func getMacroLogId(w http.ResponseWriter, r *http.Request) {
	macrolog.GetMacroLogId(w, r)
}

func deleteMacroLog(w http.ResponseWriter, r *http.Request) {
	macrolog.DeleteMacroLog(w, r)
}