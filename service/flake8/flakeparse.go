// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package flake8

import (
	"bufio"
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"os"
	"regexp"
	"strconv"
)

type MessageLine struct {
	File      string
	Line      int
	Character int
	ErrorCode string
	Message   string
	Raw       string
}

// https://regex101.com/r/yzGkd3/1
// E***/W***: pep8 errors and warnings
// F***: PyFlakes codes (see below)
// C9**: McCabe complexity plugin mccabe
// N8**: Naming Conventions plugin pep8-naming
var regExFlakeLine = regexp.MustCompile(`(?mU)\.\/(?P<File>.*):(?P<Line>\d*):(?P<Character>\d*): (?P<ErrorCode>[EWFNC]\d*) (?P<Message>.*)$`)

func MessagesToString(msgs []*MessageLine) string {
	ret := "```\n"
	for _, msg := range msgs {
		ret = ret + "\n" + msg.Raw
	}
	ret = ret + "\n```"
	return ret
}

func Parse(fn string) []*MessageLine {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var MessageLines []*MessageLine

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		match := regExFlakeLine.FindStringSubmatch(line)

		paramsMap := &MessageLine{}
		paramsMap.Raw = line
		for i, name := range regExFlakeLine.SubexpNames() {
			switch name {
			case "File":
				paramsMap.File = match[i]
				break
			case "Line":
				paramsMap.Line, _ = strconv.Atoi(match[i])
				break
			case "Character":
				paramsMap.Character, _ = strconv.Atoi(match[i])
				break
			case "ErrorCode":
				paramsMap.ErrorCode = match[i]
				break
			case "Message":
				paramsMap.Message = match[i]
			default:
				break
			}
		}

		MessageLines = append(MessageLines, paramsMap)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return MessageLines
}

func BuildReport(msgs []*MessageLine) *github.CheckRunOutput {

	var annotations []*github.CheckRunAnnotation

	countWarnings := 0
	countFailures := 0

	for _, msg := range msgs {
		annotations = append(annotations, msg.ToRunAnnotation())

		if msg.WarningLevel() == "failure" {
			countFailures++
		}
		if msg.WarningLevel() == "warning" {
			countWarnings++
		}
	}

	return &github.CheckRunOutput{
		Title:       github.String("Flake8-Check"),
		Summary:     github.String(fmt.Sprintf("There are %d failures and %d warnings.", countFailures, countWarnings)),
		Text:        github.String(MessagesToString(msgs)),
		Annotations: annotations,
	}

}

func (msg *MessageLine) WarningLevel() string {
	// "failure, warning, notice"
	if string(msg.ErrorCode[0]) == "W" {
		return "warning"
	} else {
		return "failure"
	}
}

func (msg *MessageLine) ToRunAnnotation() *github.CheckRunAnnotation {
	return &github.CheckRunAnnotation{
		FileName:     github.String(msg.File),
		BlobHRef:     github.String("https://example.com"),
		StartLine:    github.Int(msg.Line),
		EndLine:      github.Int(msg.Line),
		WarningLevel: github.String(msg.WarningLevel()),
		Message:      github.String(msg.Message),
		RawDetails:   github.String(msg.Raw),
		Title:        github.String("Flake8-Check"),
	}
}
