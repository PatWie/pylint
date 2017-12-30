package pylint

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goji.io/pat"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	// "github.com/patwie/pylint/pylint"
)

// Make an enqueuer with a particular namespace
// var enqueuer = work.NewEnqueuer("pylint_go", pylint.RedisPool)

func writeResponse(rw http.ResponseWriter, msg_text string) {
	msg := HookResponse{msg_text}
	js, err := json.Marshal(msg)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

func HandleHome(w http.ResponseWriter, req *http.Request) {
	writeResponse(w, "PyLint Go is successfully running")

}

func HandleStatus(w http.ResponseWriter, req *http.Request) {
	writeResponse(w, "org "+pat.Param(req, "org")+
		"name "+pat.Param(req, "name"))

}

func HandleReports(rw http.ResponseWriter, req *http.Request) {
	commit := pat.Param(req, "commit")

	match, _ := regexp.MatchString("([a-f0-9]{40})", commit)
	if !match {
		http.Error(rw, "400 Bad Request - Not a valid checksum", http.StatusForbidden)
		return
	}

	body, err := ioutil.ReadFile("/data/reports/" + commit)

	if err != nil {
		http.Error(rw, "404 Bad Request - Report not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(rw, "%s", body)

}

func HandleHooks(w http.ResponseWriter, r *http.Request) {
	log.Println("webhook triggered")

	// verify signature
	signature := r.Header.Get("X-Hub-Signature")
	if len(signature) == 0 {
		http.Error(w, "403 Forbidden - Missing X-Hub-Signature required for HMAC verification", http.StatusForbidden)
		return
	}
	payload, _ := ioutil.ReadAll(r.Body)

	actualMAC := string(signature[len("sha1="):])
	mac := hmac.New(sha1.New, []byte(Cfg.Github.Secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(actualMAC), []byte(expectedMAC)) {
		http.Error(w, "403 Forbidden - HMAC verification failed", http.StatusForbidden)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	if len(event) == 0 {
		http.Error(w, "400 Bad Request - Missing X-GitHub-Event Header", http.StatusBadRequest)
		return
	}
	gitHubEvent := string(event)
	log.Println("received event: " + gitHubEvent)

	switch gitHubEvent {
	case "ping":
		writeResponse(w, "pong")

	case "integration_installation":
		InstallEvent(w, payload)

	case "installation":
		InstallEvent(w, payload)

	case "push":
		PushEvent(w, payload)
	}

}

func PushEvent(w http.ResponseWriter, payload []byte) {
	var cc PushPayload
	var installation DBInstallation

	json.Unmarshal([]byte(payload), &cc)

	query := *Database.First(&installation, "installation = ?", cc.Installation.ID)
	if query.RecordNotFound() == true {
		http.Error(w, "403 Forbidden - Installation is unkown", http.StatusForbidden)
		return
	}

	// branch name is unfortunately not directly in the payload
	repo_branches := strings.Split(cc.Ref, "/")
	repo_branch := repo_branches[len(repo_branches)-1]

	_, err := enqueue(Cfg.Github.IntegrationID,
		cc.Installation.ID,
		cc.After,
		cc.Repository.Owner.Name,
		cc.Repository.Name,
		repo_branch)

	if err != nil {
		log.Fatal(err)
	}
}

func InstallEvent(w http.ResponseWriter, payload []byte) {
	var cc InstallationPayload
	json.Unmarshal([]byte(payload), &cc)

	if Cfg.Github.AdminId != 0 {
		if cc.Sender.ID != Cfg.Github.AdminId {
			http.Error(w, "403 Forbidden - User Id verification failed", http.StatusForbidden)
			return
		}

	}

	switch cc.Action {
	case "created":
		CreateInstallation(cc)
		writeResponse(w, "integration is successfully installed")
	case "deleted":
		DeleteInstallation(cc)
		writeResponse(w, "integration is successfully uninstalled")
	}
}
