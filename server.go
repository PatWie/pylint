package main

// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// wrapper for database
var db *gorm.DB
var cfg PyLintConfig

// Make an enqueuer with a particular namespace
var enqueuer = work.NewEnqueuer("pylint_go", RedisPool)

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

// Serve flake8 reports.
func HandleHome(w http.ResponseWriter, req *http.Request) {
	writeResponse(w, "PyLint Go is successfully running")

}

func HandleReports(rw http.ResponseWriter, req *http.Request) {
	commit := req.URL.Path[len("/report/"):]

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
	actualMAC := string(signature[len("sha1="):])

	payload, _ := ioutil.ReadAll(r.Body)

	mac := hmac.New(sha1.New, []byte(cfg.Github.Secret))
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
		var cc InstallationPayload
		json.Unmarshal([]byte(payload), &cc)
		if cfg.Github.AdminId != 0 {
			if cc.Sender.ID != cfg.Github.AdminId {
				http.Error(w, "403 Forbidden - User Id verification failed", http.StatusForbidden)
				return
			}

		}
		if cc.Action == "created" {
			db.Create(&DBInstallation{Sender: cc.Sender.ID,
				Installation: cc.Installation.ID})
			writeResponse(w, "integration is successfully installed")
		}
		if cc.Action == "deleted" {
			db.Where("installation = ?", cc.Installation.ID).Delete(DBInstallation{})
			writeResponse(w, "integration is successfully uninstalled")
		}

	case "installation":
		// add this app to repository
		var cc InstallationPayload
		json.Unmarshal([]byte(payload), &cc)
		if cfg.Github.AdminId != 0 {
			if cc.Sender.ID != cfg.Github.AdminId {
				http.Error(w, "403 Forbidden - User Id verification failed", http.StatusForbidden)
				return
			}
		}

		if cc.Action == "created" {
			db.Create(&DBInstallation{Sender: cc.Sender.ID,
				Installation: cc.Installation.ID})
			writeResponse(w, "integration is successfully installed")
		}
		if cc.Action == "deleted" {
			db.Where("installation = ?", cc.Installation.ID).Delete(DBInstallation{})
			writeResponse(w, "integration is successfully uninstalled")
		}

	case "push":
		var cc PushPayload
		json.Unmarshal([]byte(payload), &cc)
		log.Println(cc.Installation.ID)
		log.Println(cc.Sender.ID)

		var installation DBInstallation
		query := *db.First(&installation, "installation = ?", cc.Installation.ID)
		if query.RecordNotFound() == true {
			http.Error(w, "403 Forbidden - Installation is unkown", http.StatusForbidden)
			return
		}

		// Use installation transport with client.
		log.Println("Update commit " + cc.After)
		log.Println("github_integrationID " + strconv.FormatInt(cfg.Github.IntegrationID, 10))
		log.Println("cc.Installation.ID " + strconv.FormatInt(cc.Installation.ID, 10))

		_, err := enqueuer.Enqueue("test_repo",
			work.Q{
				"integration_id":  strconv.FormatInt(cfg.Github.IntegrationID, 10),
				"installation_id": strconv.FormatInt(cc.Installation.ID, 10),
				"commit_sha1":     cc.After,
				"repo_owner":      cc.Repository.Owner.Name,
				"repo_name":       cc.Repository.Name})
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Unable to parse config: ", err)
	}
	if err := env.Parse(&cfg.Github); err != nil {
		log.Fatal("Unable to parse config.GitHub: ", err)
	}
	if err := env.Parse(&cfg.Database); err != nil {
		log.Fatal("Unable to parse config.Database: ", err)
	}

	// database
	// -------------------------------------------------------
	var err error
	db, err = gorm.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&DBInstallation{})

	// http server
	// -------------------------------------------------------
	log.Println("start application and listen on (internal):", cfg.Port)
	log.Println("start application and listen on (public):", cfg.PublicPort)
	http.HandleFunc("/report/", HandleReports)
	http.HandleFunc("/hook", HandleHooks)
	http.HandleFunc("/home", HandleHome)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)

}
