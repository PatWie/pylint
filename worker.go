package main

import (
	"bufio"
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
	"os/exec"
	"os/signal"
	"strconv"
)

type Workload struct {
	InstallationId int
	Commit         struct {
		Organization string
		Name         string
		Branch       string
		Sha1         string
	}
	Result struct {
		Status string
	}
}

func NewWorkload(id string, org string, name string, branch string, sha1 string) Workload {
	w := Workload{}
	i, _ := strconv.Atoi(id)
	w.InstallationId = i
	w.Commit.Organization = org
	w.Commit.Name = name
	w.Commit.Branch = branch
	w.Commit.Sha1 = sha1

	w.Result.Status = "pending"
	return w
}

func (wl *Workload) SendStatus(ctx *context.Context, client *github.Client, status string) (*github.RepoStatus, error) {

	state_url := pylint.Cfg.Url + "report/" + wl.Commit.Sha1
	if status == "pending" {
		state_url = pylint.Cfg.Url + ""
	}

	// convert to struct
	new_state := github.RepoStatus{
		State:       &status,
		TargetURL:   &state_url,
		Description: &pylint.Cfg.Name}

	// start sending
	statuses, _, err := client.Repositories.CreateStatus(
		*ctx,
		wl.Commit.Organization,
		wl.Commit.Name,
		wl.Commit.Sha1,
		&new_state)
	return statuses, err

}

func (wl *Workload) RunTest(token string) ([]byte, error) {
	cmd, err := exec.Command("/go/src/pylint/Docker/run.sh",
		wl.Commit.Sha1,
		token,
		wl.Commit.Organization,
		wl.Commit.Name).Output()
	fmt.Println("done check")
	fmt.Println(string(cmd))
	return cmd, err
}

func (wl *Workload) LintStatus() pylint.DBLintStatus {
	return pylint.DBLintStatus{Organization: wl.Commit.Organization,
		Repository: wl.Commit.Name,
		Branch:     wl.Commit.Branch}
}

func (wl *Workload) GetResult() (int, error) {
	file, err := os.Open("/data/reports/" + wl.Commit.Sha1)
	if err != nil {
		return 0, err
		// _, _ = sendCommitStatus(ctx, commit_sha1,
		// 	repo_owner, repo_name, pylint.GIT_STATUS_FAILURE, client)
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount, nil
}

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

	pylint.Database.AutoMigrate(&pylint.DBInstallation{})
	pylint.Database.AutoMigrate(&pylint.DBLintStatus{})

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

	wl := NewWorkload(
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
	if terr != nil {
		log.Println("cannot fetch access_token: " + terr.Error())
		return nil
	}

	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()

	wl.SendStatus(&ctx, client, "pending")

	// -------------------------------------------------
	_, err = wl.RunTest(access_token)
	if err != nil {
		log.Println("cannot run pylintint script: " + err.Error())
		return nil
	}

	lineCount, err := wl.GetResult()
	if err != nil || lineCount > 0 {
		wl.SendStatus(&ctx, client, "failed")
	} else {
		wl.SendStatus(&ctx, client, "success")
	}

	lintstatus := pylint.DBLintStatus{}
	pylint.Database.Where(wl.LintStatus()).FirstOrInit(&lintstatus)
	lintstatus.Status = lineCount
	pylint.Database.Save(&lintstatus)

	fmt.Println("number of lines from flake8:", lineCount)

	return nil
}
