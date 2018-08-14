// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package service

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/google/go-github/github"
	"github.com/patwie/pylint/model"
	"github.com/patwie/pylint/service/flake8"
	"log"
	"net/http"
	"time"
)

const (
	GIT_STATUS_FAILURE string = "failure"
	GIT_STATUS_PENDING string = "pending"
	GIT_STATUS_SUCCESS string = "success"
)

var config = model.GetConfiguration()

// Make a redis pool
var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		// return redis.Dial("tcp", "redis:6379")
		return redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port))
	},
}

type WorkerContext struct{}

func (c *WorkerContext) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Println("Starting job: ", job.Name)
	return next()
}

func (c *WorkerContext) TestRepo(job *work.Job) error {
	log.Println("TestRepo started:", job)

	installationID := job.ArgInt64("installation_id")
	commitSHA := job.ArgString("sha")
	repoOwner := job.ArgString("owner")
	repoName := job.ArgString("repository")
	repoBranch := job.ArgString("branch")

	log.Println("IntegrationID  ", config.GitHub.IntegrationID)
	log.Println("installationID  ", installationID)
	log.Println("commitSHA  ", commitSHA)
	log.Println("repoOwner  ", repoOwner)
	log.Println("repoName  ", repoName)
	log.Println("repoBranch  ", repoBranch)

	if err := job.ArgError(); err != nil {
		log.Println(err)
		return err
	}

	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		int(config.GitHub.IntegrationID),
		int(installationID),
		config.Pylint.KeyFile)

	log.Printf("run for installation %v", installationID)

	if err != nil {
		log.Println("cannot generate token: " + err.Error())
		return nil
	} else {
		log.Println("got token")
	}

	access_token, terr := itr.Token()

	if terr != nil {
		log.Println("cannot fetch access_token: " + terr.Error())
		return nil
	} else {
		log.Println("fetched access_token")
		log.Println(access_token)
	}

	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()

	startTime := time.Now()
	log.Println("started at", startTime)

	// set pending
	opt := github.CreateCheckRunOptions{
		Name:       config.Pylint.Name,
		HeadBranch: repoBranch,
		HeadSHA:    commitSHA,
		DetailsURL: github.String(fmt.Sprintf("%s/", config.Pylint.URL)),
		StartedAt:  &github.Timestamp{startTime},
		Status:     github.String("in_progress"),
		Output: &github.CheckRunOutput{
			Title:   github.String("Test report"),
			Summary: github.String("Started to run flake8"),
			Text:    github.String(""),
		},
	}

	log.Println(opt)
	_, _, err = client.Checks.CreateCheckRun(ctx, repoOwner, repoName, opt)
	if err != nil {
		log.Println(err)
		return nil
	}

	messages, err := RunLinter(access_token, repoOwner, repoName, commitSHA)

	log.Println("linter finished")

	opt = github.CreateCheckRunOptions{
		Name:        config.Pylint.Name,
		HeadBranch:  repoBranch,
		HeadSHA:     commitSHA,
		DetailsURL:  github.String(fmt.Sprintf("%s/%s/%s/report/%s", config.Pylint.URL, repoOwner, repoName, commitSHA)),
		Status:      github.String("completed"),
		StartedAt:   &github.Timestamp{startTime},
		CompletedAt: &github.Timestamp{time.Now()},
		Output:      flake8.BuildReport(messages),
	}

	if len(messages) > 0 {
		opt.Conclusion = github.String("failure")
	} else {
		opt.Conclusion = github.String("success")
	}

	_, _, err = client.Checks.CreateCheckRun(ctx, repoOwner, repoName, opt)
	if err != nil {
		panic(err)
	}

	return nil
}

func TestRedisActive() {
	conn := RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		log.Println(err)
		panic("Can't connect to the Redis database")
	}
}

var enqueuer = createEnqueuer()
var dequeuer = createDequeuer()

func createEnqueuer() *work.Enqueuer {
	return work.NewEnqueuer("python_linter_pool", RedisPool)
}

func createDequeuer() *work.WorkerPool {
	pool := work.NewWorkerPool(WorkerContext{}, 10, "python_linter_pool", RedisPool)
	pool.Middleware((*WorkerContext).Log)
	pool.Job("test_repo", (*WorkerContext).TestRepo)
	return pool
}

func GetEnqueuer() *work.Enqueuer {
	return enqueuer
}

func GetDequeuer() *work.WorkerPool {
	return dequeuer
}
