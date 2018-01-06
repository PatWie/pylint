package pylint

import (
	"bufio"
	"context"
	"fmt"

	"github.com/google/go-github/github"

	"io/ioutil"
	"os"
	"os/exec"

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
		Report string
		Lines  int
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
	w.Result.Status = GIT_STATUS_PENDING
	w.Result.Report = ""
	w.Result.Lines = 0
	return w
}

func (wl *Workload) SendStatus(ctx *context.Context, client *github.Client, status string) (*github.RepoStatus, error) {

	state_url := ""
	if status == GIT_STATUS_PENDING {
		state_url = Cfg.Url
	} else {
		state_url = fmt.Sprintf("%s/%s/%s/%s/report",
			Cfg.Url,
			wl.Commit.Organization,
			wl.Commit.Name,
			wl.Commit.Sha1)
	}

	// convert to struct
	new_state := github.RepoStatus{
		State:       &status,
		TargetURL:   &state_url,
		Description: &Cfg.Name}

	// start sending
	statuses, _, err := client.Repositories.CreateStatus(
		*ctx,
		wl.Commit.Organization,
		wl.Commit.Name,
		wl.Commit.Sha1,
		&new_state)
	return statuses, err

}

func (wl *Workload) RunTest(token string) (string, error) {
	cmd, err := exec.Command("/go/src/github.com/patwie/pylint/Docker/run.sh",
		wl.Commit.Sha1,
		token,
		wl.Commit.Organization,
		wl.Commit.Name).Output()
	if err != nil {
		fmt.Println("bash: " + err.Error())
		return "", err
	} else {
		fmt.Println("done check")
		fmt.Println(string(cmd))
		return string(cmd), nil

	}
}

func (wl *Workload) LintStatus() DBLintStatus {
	return DBLintStatus{Organization: wl.Commit.Organization,
		Repository: wl.Commit.Name,
		Branch:     wl.Commit.Branch}
}

func (wl *Workload) GetResult() error {
	reportsDir := "/data/reports/"

	file, err := os.Open(reportsDir + wl.Commit.Sha1)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}

	wl.Result.Lines = lineCount

	b, err := ioutil.ReadFile(reportsDir + wl.Commit.Sha1)
	if err != nil {
		return err
	}
	wl.Result.Report = string(b)
	return nil
}
