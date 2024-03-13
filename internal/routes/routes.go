package routes

import (
	"net/http"

	dailyMacroTotal "github.com/ramenNoodles1998/macros-backend/internal/daily-macro-total"
	fooditem "github.com/ramenNoodles1998/macros-backend/internal/food-item"
	macrolog "github.com/ramenNoodles1998/macros-backend/internal/macro-log"
)
func Router() http.Handler {
	mux := http.NewServeMux()

	macrolog.SetMacroLogRoutes(mux)
	fooditem.SetFoodItemRoutes(mux)
	dailyMacroTotal.SetMacroLogRoutes(mux)

	return mux
}

