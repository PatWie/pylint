// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package main

import (
	"fmt"
	"github.com/patwie/pylint/model"
	"github.com/patwie/pylint/router"
	"github.com/patwie/pylint/service"
	"net/http"
)

func main() {
	service.TestRedisActive()
	config := model.GetConfiguration()

	r := router.GetRouter()

	fmt.Printf("ListenAndServe at :%d", config.Pylint.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Pylint.Port), r)
	if err != nil {
		panic(err)
	}
}
