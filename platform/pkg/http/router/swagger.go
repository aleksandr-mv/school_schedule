package http_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MountSwagger(r *chi.Mux, swaggerPath, uiPath string) {
	r.Get("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, uiPath)
	})

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
	})
}
