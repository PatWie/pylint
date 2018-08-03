// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package api

import (
	"github.com/patwie/pylint/router/render"
	"net/http"
)

// Default homepage.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	render.WriteText(w, "active")
}

// Return version as JSON.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	render.WriteJSON(w, render.H{
		"source":  "https://github.com/patwie/pylint",
		"version": "0.0.1"},
	)
}

// Return pong to ping.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	render.WriteText(w, "pong")
}
