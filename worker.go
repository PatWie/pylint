package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/caarlos0/env"
	// "github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
)

var cfg PyLintConfig

// Send new commit status to GitHub server.
func sendCommitStatus(
	ctx context.Context,
	commit string,
	user_name string,
	repo_name string,
	status string,
	client *github.Client) (*github.RepoStatus, error) {

	state_url := cfg.Url + "report/" + commit
	if status == "pending" {
		state_url = cfg.Url + ""
	}

	// convert to struct
	new_state := github.RepoStatus{State: &status,
		TargetURL:   &state_url,
		Description: &cfg.Name}

	// start sending
	statuses, _, err := client.Repositories.CreateStatus(ctx, user_name,
		repo_name, commit, &new_state)
	return statuses, err
}

var db *gorm.DB

type WContext struct{}

func main() {
	var err error

	// test redis
	conn := RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("PING")
	if err != nil {
		log.Fatal("Can't connect to the Redis database")
	}

	// read config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Unable to parse config: ", err)
	}
	if err := env.Parse(&cfg.Github); err != nil {
		log.Fatal("Unable to parse config.GitHub: ", err)
	}
	if err := env.Parse(&cfg.Database); err != nil {
		log.Fatal("Unable to parse config.Database: ", err)
	}

	log.Println("worker is ready ...")

	// database
	db, err = gorm.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&DBInstallation{})

	pool := work.NewWorkerPool(WContext{}, 10, "pylint_go", RedisPool)
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
	installation_id := job.ArgString("installation_id")
	commit_sha1 := job.ArgString("commit_sha1")
	repo_owner := job.ArgString("repo_owner")
	repo_name := job.ArgString("repo_name")

	if err := job.ArgError(); err != nil {
		log.Println(err)
		return err
	}

	log.Println("worker called")
	log.Println("installation_id " + installation_id)
	log.Println("commit_sha1 " + commit_sha1)
	log.Println("repo_owner " + repo_owner)
	log.Println("repo_name " + repo_name)

	installation_id_int, err := strconv.Atoi(installation_id)
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport,
		int(cfg.Github.IntegrationID),
		installation_id_int,
		cfg.Github.KeyPath)

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
	_, _ = sendCommitStatus(ctx, commit_sha1,
		repo_owner, repo_name, GIT_STATUS_PENDING, client)

	// -------------------------------------------------

	cmd, err := exec.Command("/go/src/pylint/run.sh", commit_sha1, access_token, repo_owner, repo_name).Output()
	if err != nil {
		log.Println("cannot run pylintint script: " + err.Error())
		return nil
	}

	fmt.Println("done check")
	fmt.Println(string(cmd))

	file, err := os.Open("/data/reports/" + commit_sha1)
	if err != nil {
		_, _ = sendCommitStatus(ctx, commit_sha1,
			repo_owner, repo_name, GIT_STATUS_FAILURE, client)
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}

	if lineCount > 0 {
		_, _ = sendCommitStatus(ctx, commit_sha1,
			repo_owner, repo_name, GIT_STATUS_FAILURE, client)
	} else {
		_, _ = sendCommitStatus(ctx, commit_sha1,
			repo_owner, repo_name, GIT_STATUS_SUCCESS, client)
	}

	fmt.Println("number of lines from flake8:", lineCount)

	return nil
}
