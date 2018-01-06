package pylint

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"strconv"
)

// Make an enqueuer with a particular namespace
var Enqueuer *work.Enqueuer

func CreateQueue(pool *redis.Pool) {
	Enqueuer = work.NewEnqueuer("pylint_go", RedisPool)
}

func enqueue(integration_id int64,
	installation_id int64,
	sha1 string,
	owner string,
	repo string,
	branch string) (*work.Job, error) {
	return Enqueuer.Enqueue("test_repo",
		work.Q{
			"integration_id":  strconv.FormatInt(integration_id, 10),
			"installation_id": strconv.FormatInt(installation_id, 10),
			"commit_sha1":     sha1,
			"repo_owner":      owner,
			"repo_name":       repo,
			"repo_branch":     branch})
}
