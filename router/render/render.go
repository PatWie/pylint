// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package render

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// similar to gin.H as a neat wrapper

type H map[string]interface{}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, []string{"application/json; charset=utf-8"})
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func WriteText(w http.ResponseWriter, format string) error {
	writeContentType(w, []string{"text/plain; charset=utf-8"})
	io.WriteString(w, format)
	return nil
}

func WriteTextf(w http.ResponseWriter, format string, a ...interface{}) error {
	writeContentType(w, []string{"text/plain; charset=utf-8"})
	io.WriteString(w, fmt.Sprintf(format, a...))
	return nil
}

// Example:
// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 	render.WriteJSON(w, render.H{"test": "hi"})
// 	render.WriteText(w, "hi")
// 	http.Redirect(w, r, "/", 301)
// })
