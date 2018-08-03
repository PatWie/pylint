// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package main

import (
	"github.com/patwie/pylint/model"
	"github.com/patwie/pylint/service"
	"log"
	"os"
	"os/signal"
)

func main() {
	config := model.GetConfiguration()
	config.Debug()

	log.Printf("Start worker for integration %v", config.GitHub.IntegrationID)
	service.TestRedisActive()

	pool := service.GetDequeuer()
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}
