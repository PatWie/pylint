// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/patwie/pylint/router/api"
)

func GetRouter() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)

	r.Get("/", api.HomeHandler)

	r.Route("/{user}", func(r chi.Router) {
		r.Route("/{repo}", func(r chi.Router) {
			r.Route("/report", func(r chi.Router) {
				r.Get("/{commit}", api.ReportHandler)
			})
		})
	})

	r.Get("/version", api.VersionHandler)
	r.Get("/ping", api.PingHandler)
	r.Post("/hook", api.HookHandler)

	return r

}
