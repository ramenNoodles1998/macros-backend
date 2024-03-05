package routes

import (
	"net/http"

	fooditem "github.com/ramenNoodles1998/macros-backend/internal/food-item"
	macrolog "github.com/ramenNoodles1998/macros-backend/internal/macro-log"
)
func Router() http.Handler {
	mux := http.NewServeMux()

	macrolog.SetMacroLogRoutes(mux)
	fooditem.SetFoodItemRoutes(mux)

	return mux
}

