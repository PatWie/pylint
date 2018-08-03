// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/patwie/pylint/router/api"
	"github.com/patwie/pylint/router/render"
	"net/http"
)

func GetRouter() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.WriteText(w, "active")
	})

	r.Route("/{user}", func(r chi.Router) {
		r.Route("/{repo}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				render.WriteJSON(w,
					render.H{
						"user": chi.URLParam(r, "user"),
						"repo": chi.URLParam(r, "repo")},
				)
			})
			r.Route("/report", func(r chi.Router) {
				r.Get("/{commit}", api.ReportHandler)
			})
		})
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.WriteJSON(w, render.H{
			"source":  "https://github.com/patwie/pylint",
			"version": "0.0.1"},
		)
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.WriteText(w, "pong")
	})

	r.Post("/hook", api.HookHandler)

	return r

}
