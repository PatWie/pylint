package main

import (
	// "bufio"
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation"
	// "github.com/caarlos0/env"
	// "github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/patwie/pylint/pylint"
	"log"
	"net/http"
	"os"
	// "os/exec"
	"os/signal"
	// "strconv"
)

type WContext struct{}

var db *gorm.DB

func main() {
	var err error

	// read config
	pylint.Cfg.Parse()

	// test redis
	redis := pylint.ConnectRedis(pylint.Cfg)
	defer redis.Close()
	_, err = redis.Do("PING")
	if err != nil {
		log.Println(err)
		log.Fatal("Can't connect to the Redis database")
	}

	log.Println("worker is ready ...")

	err = pylint.ConnectDatabase(pylint.Cfg)
	if err != nil {
		panic("failed to connect database")
	}
	defer pylint.Database.Close()
	pylint.MigrateDatabase(pylint.Database)

	pool := work.NewWorkerPool(WContext{}, 10, "pylint_go", pylint.RedisPool)
	pool.Middleware((*WContext).Log)
	pool.Job("test_repo", (*WContext).TestRepo)
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *WContext) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *WContext) TestRepo(job *work.Job) error {

	wl := pylint.NewWorkload(
		job.ArgString("installation_id"),
		job.ArgString("repo_owner"),
		job.ArgString("repo_name"),
		job.ArgString("repo_branch"),
		job.ArgString("commit_sha1"))

	if err := job.ArgError(); err != nil {
		log.Println(err)
		return err
	}

	log.Println("worker called")

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport,
		int(pylint.Cfg.Github.IntegrationID),
		wl.InstallationId,
		pylint.Cfg.Github.KeyPath)
	if err != nil {
		log.Println("cannot generate token: " + err.Error())
		return nil
	}

	access_token, terr := itr.Token()
	log.Println(access_token)
	if terr != nil {
		log.Println("cannot fetch access_token: " + terr.Error())
		return nil
	}

	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()

	wl.SendStatus(&ctx, client, pylint.GIT_STATUS_PENDING)

	// -------------------------------------------------
	cmd, err := wl.RunTest(access_token)
	log.Println(cmd)
	log.Println(err.Error())
	// // if err != nil {
	// // 	log.Println("cannot run pylintint script: " + err.Error())
	// // 	return nil
	// // }

	// log.Println("alive")
	// err = wl.GetResult()
	// if err != nil || wl.Result.Lines > 0 {
	// 	log.Println(err.Error())
	// 	wl.SendStatus(&ctx, client, pylint.GIT_STATUS_FAILURE)
	// } else {
	// 	wl.SendStatus(&ctx, client, pylint.GIT_STATUS_SUCCESS)
	// }
	// log.Println("alive2")

	// // lintstatus := pylint.DBLintStatus{}
	// // pylint.Database.Where(wl.LintStatus()).FirstOrInit(&lintstatus)
	// // lintstatus.Status = wl.Result.Lines
	// // pylint.Database.Save(&lintstatus)

	fmt.Println("number of lines from flake8:", wl.Result.Lines)

	return nil
}
