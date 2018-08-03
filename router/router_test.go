// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package router

import (
	"encoding/json"
	. "github.com/franela/goblin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {

	r := GetRouter()

	g := Goblin(t)
	g.Describe("Router", func() {

		g.It("Should return the version", func() {

			expectedBody, _ := json.Marshal(map[string]string{
				"source":  "https://github.com/patwie/pylint",
				"version": "0.0.1",
			})

			req, _ := http.NewRequest("GET", "/version", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			g.Assert(resp.Code).Equal(http.StatusOK)
			g.Assert(resp.Body.String()).Equal(string(expectedBody))
		})

		g.It("Should play ping pong", func() {
			req, _ := http.NewRequest("GET", "/ping", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			g.Assert(resp.Code).Equal(http.StatusOK)
			g.Assert(resp.Body.String()).Equal(string("pong"))
		})

	})

}
