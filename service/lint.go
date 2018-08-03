// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package service

import (
	"bufio"
	"fmt"
	"github.com/patwie/pylint/model"
	"github.com/patwie/pylint/service/flake8"
	"log"
	"os"
	"os/exec"
)

type LintError struct {
	lineCount int
}

func (e LintError) Error() string {
	return fmt.Sprintf("Linter found %i errors", e.lineCount)
}

func RunLinter(access_token string, repoOwner string, repoName string, commitSHA string) ([]*flake8.MessageLine, error) {

	config := model.GetConfiguration()

	// run linter
	ss, err := exec.Command("./scripts/run_job.sh", commitSHA, access_token, repoOwner, repoName, config.Pylint.ReportsPath).Output()
	if err != nil {
		log.Println("cannot run pylinting script: " + err.Error())
		return nil, err
	}
	fmt.Println(string(ss))

	// scan report
	file, err := os.Open(config.Pylint.ReportsPath + "/" + commitSHA)
	if err != nil {
		log.Fatalf("did not find file %s", config.Pylint.ReportsPath+"/"+commitSHA)
		return nil, err
	}

	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}

	messages := flake8.Parse(config.Pylint.ReportsPath + "/" + commitSHA)
	return messages, nil
}
