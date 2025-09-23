package http_router

import (
	"github.com/go-chi/chi/v5"

	healthhttp "github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/http/health"
)

// RegisterHealthRoutes регистрирует стандартные health check маршруты
func RegisterHealthRoutes(r chi.Router) {
	r.Get("/healthz", healthhttp.Handler)
	r.Get("/live", healthhttp.LiveHandler)
	r.Get("/ready", healthhttp.ReadyHandler)
	r.Get("/start", healthhttp.StartHandler)
}
