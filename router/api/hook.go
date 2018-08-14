// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package api

import (
	"github.com/google/go-github/github"
	"github.com/patwie/pylint/model"
	"github.com/patwie/pylint/router/render"
	"github.com/patwie/pylint/service"
	"github.com/patwie/pylint/store"
	"log"
	"net/http"
)

var enqueuer = service.GetEnqueuer()
var config = model.GetConfiguration()

// Handle all incoming GitHub hooks
func HookHandler(w http.ResponseWriter, r *http.Request) {
	config := model.GetConfiguration()

	payload, err := github.ValidatePayload(r, []byte(config.GitHub.Secret))
	if err != nil {
		http.Error(w, "400 Bad Request - Invalid request body (GitHub secret might be invalid)", http.StatusBadRequest)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		http.Error(w, "400 Bad Request - Could not parse webhook", http.StatusBadRequest)
		return
	}

	switch e := event.(type) {

	// PingEvent is triggered when a Webhook is added to GitHub.
	//
	// GitHub API docs: https://developer.github.com/webhooks/#ping-event
	case *github.PingEvent:
		render.WriteText(w, "pong")
		return

	// InstallationEvent is triggered when a GitHub App has been installed or uninstalled.
	// The Webhook event name is "installation".
	//
	// GitHub API docs: https://developer.github.com/v3/activity/events/types/#installationevent
	case *github.InstallationEvent:
		if config.Pylint.AdminId != 0 {
			// we restrict installations only to installations from the specified user
			if config.Pylint.AdminId != *e.GetSender().ID {
				http.Error(w, "400 Bad Request - Installation API is restricted", http.StatusBadRequest)
			}
		}

		switch *e.Action {
		case "created":
			store.DS().CreateInstallation(&model.Installation{
				SenderID:       int(*e.GetSender().ID),
				InstallationID: int(*e.GetInstallation().ID),
			})

		case "deleted":
			store.DS().DeleteInstallation(int(*e.GetInstallation().ID))
		}

	// CheckRunEvent is triggered when a check run is "created", "updated", or "re-requested".
	// The Webhook event name is "check_run".
	//
	// GitHub API docs: https://developer.github.com/v3/activity/events/types/#checkrunevent
	case *github.CheckSuiteEvent:
		b := &model.Build{
			InstallationID: int(*e.Installation.ID),
			SHA:            *e.CheckSuite.HeadSHA,
			Owner:          *e.Repo.Owner.Login,
			Repository:     *e.Repo.Name,
			Branch:         *e.CheckSuite.HeadBranch,
		}

		_, err := enqueuer.Enqueue("test_repo", b.Serialize())
		if err != nil {
			http.Error(w, "500 Internal Server Error - See logs", http.StatusInternalServerError)
			log.Fatalln(err)
		}

	// CheckSuiteEvent is triggered when a check suite is "completed", "requested", or "re-requested".
	// The Webhook event name is "check_suite".
	//
	// GitHub API docs: https://developer.github.com/v3/activity/events/types/#checksuiteevent
	case *github.CheckRunEvent:
		if *e.Action == "requested" || *e.Action == "re-requested" {
			b := &model.Build{
				InstallationID: int(*e.Installation.ID),
				SHA:            *e.CheckRun.CheckSuite.HeadSHA,
				Owner:          *e.Repo.Owner.Login,
				Repository:     *e.Repo.Name,
				Branch:         *e.CheckRun.CheckSuite.HeadBranch,
			}

			_, err := enqueuer.Enqueue("test_repo", b.Serialize())
			if err != nil {
				http.Error(w, "500 Internal Server Error - See logs", http.StatusInternalServerError)
				log.Fatalln(err)
			}
		}

	default:
		// log.Printf("unknown event type %s\n", github.WebHookType(r))
		render.WriteTextf(w, "Request ok - but I do not speak Klingon - unknown event type %s", github.WebHookType(r))
		return

	}
	render.WriteText(w, "done")
}
