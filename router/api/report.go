// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package api

import (
	"github.com/go-chi/chi"
	"github.com/patwie/pylint/model"
	"github.com/patwie/pylint/router/render"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Return flake8 text files
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	config := model.GetConfiguration()
	commit := chi.URLParam(r, "commit")

	match, _ := regexp.MatchString("([a-f0-9]{40})", commit)
	if !match {
		http.Error(w, "400 Bad Request - Not a valid checksum", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadFile(config.Pylint.ReportsPath + "/" + commit)

	if err != nil {
		http.Error(w, "404 Bad Request - Report not found", http.StatusNotFound)
		return
	}
	render.WriteText(w, string(body))
}
