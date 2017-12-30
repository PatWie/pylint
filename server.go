package main

// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

import (
	"fmt"

	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/patwie/pylint/pylint"
)

func main() {

	pylint.Cfg.Parse()

	err := pylint.ConnectDatabase(pylint.Cfg)
	if err != nil {
		panic("failed to connect database")
	}
	defer pylint.Database.Close()

	redis := pylint.ConnectRedis(pylint.Cfg)
	defer redis.Close()
	_, err = redis.Do("PING")
	if err != nil {
		log.Println(err)
		log.Fatal("Can't connect to the Redis database")
	}

	pylint.CreateQueue(pylint.RedisPool)

	pylint.Database.AutoMigrate(&pylint.DBInstallation{})
	pylint.Database.AutoMigrate(&pylint.DBLintStatus{})

	// http server
	// -------------------------------------------------------
	log.Println("start application and listen on (internal):", pylint.Cfg.Port)
	log.Println("start application and listen on (public):", pylint.Cfg.PublicPort)

	root := goji.NewMux()
	root.HandleFunc(pat.Get("/hook"), pylint.HandleHooks)
	root.HandleFunc(pat.Get("/home"), pylint.HandleHome)

	repos := goji.SubMux()
	root.Handle(pat.New("/r/:org/:name/*"), repos)
	root.HandleFunc(pat.Get("/:commit"), pylint.HandleReports)
	root.HandleFunc(pat.Get("/status.svg"), pylint.HandleStatus)
	http.ListenAndServe(fmt.Sprintf(":%d", pylint.Cfg.Port), root)

}
