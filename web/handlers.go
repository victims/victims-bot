package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	githubhook "gopkg.in/rjz/githubhook.v0"
	yaml "gopkg.in/yaml.v2"

	"github.com/victims/victims-bot/cmd"
	"github.com/victims/victims-bot/gh"
	"github.com/victims/victims-bot/log"
	"github.com/victims/victims-bot/process"
)

// Hook is the main webhook endpoint
func Hook(w http.ResponseWriter, req *http.Request) {

	// Parse the hook
	hook, err := githubhook.Parse(cmd.Config.GetSecret(), req)
	if err != nil {
		log.Logger.Infof("could not read webhook: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Pass functionality to the proper hook for processing
	switch hook.Event {
	case "push":
		pushEvent(hook, w, req)
	case "ping":
		pingEvent(hook, w, req)
	default:
		log.Logger.Errorf("An unknown event was sent: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// pingEvent handles execution of github ping events
func pingEvent(hook *githubhook.Hook, w http.ResponseWriter, req *http.Request) {
	event := gh.PingEvent{}
	if err := hook.Extract(&event); err != nil {
		log.Logger.Errorf("Unable to deserialize ping event: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("Ping Event received successfully: %#v", event)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}

// pushEvent handles execution of github push events
func pushEvent(hook *githubhook.Hook, w http.ResponseWriter, req *http.Request) {
	// Grab the test variable. If it is set then we will end up
	// skipping some things so we don't end up testing external services
	// or modifying external states
	isTest := os.Getenv("VICTIMS_BOT_TEST")

	if isTest != "" {
		log.Logger.Warn("In testing mode!!!")
	}

	event := gh.PushEvent{}
	//json.Unmarshal(body, &pushEvent)
	hook.Extract(&event)
	cloneDir, err := gh.Clone(cmd.Config.GitRepo)
	if err != nil {
		log.Logger.Errorf("Web hook failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Exit out early if the push was done by the bot
	if event.Pusher.Name == cmd.Config.GitHubUsername {
		log.Logger.Infof("Ignoring push as it was done by me")
		w.WriteHeader(http.StatusOK)
		return
	}

	// The number of commits made locally
	commits := 0
	// NOTE: that if anything fails while processing the change the entire
	// set of work will not be pushed back to the repo
	for _, commit := range event.Commits {
		if commit.Author.Usersname == cmd.Config.GitHubUsername {
			log.Logger.Infof(
				"Skipping %s as it was authord by me", commit.Author.Usersname)
			continue
		}
		// Probably put this in it's own bounded goroutine
		for _, file := range commit.Added {
			log.Logger.Warnf("FILE: %s", file)
			if !strings.HasPrefix(file, "database/") {
				log.Logger.Debugf("%s isn't in the database/ path", file)
				continue
			}
			content, err := gh.GetContent(cloneDir, commit.ID, file)
			if err != nil {
				log.Logger.Errorf("Error getting contents: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO: Get the name and URL from the content
			entry := process.CVEDBEntry{}
			if err = yaml.Unmarshal([]byte(content), &entry); err != nil {
				log.Logger.Warnf("Unable to deserialize database file: %s", err)
				continue
			}
			log.Logger.Debugf("%#v", entry)
			// Continue on if the URL is not provided
			if entry.URL == "" || entry.Name == "" {
				// TODO: Maybe auto open up a GitHub issue
				log.Logger.Errorf("Unable to download package for %s due to missing data", file)
				continue
			}

			if isTest == "" {
				// TODO: Download the artifact
				// fileName, err := process.GetPackage(entry.Name, entry.URL)
				// if err != nil {
				// 	log.Logger.Errorf("Unable to download package: %s", err)
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }
				// TODO: Submit to the hash service
				// hashes, err := process.SubmitPackage(fileName, "")
				// if err := process.AddHashesToFile(file, hashes); err != nil {
				// 	log.LoggerInfof("Unable to add hash to file: %s", err)
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }
			} else {
				log.Logger.Info("Skipping package submission due to testing")
			}
			_, err = gh.CommitChange(cloneDir, file)
			if err != nil {
				log.Logger.Errorf("Error committing change: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// If we get here we've made a new local commit
			commits = commits + 1
		}
	}

	// If we have at least 1 commit then push
	if commits > 0 {
		if isTest == "" {
			// TODO: Uncomment :-)
			// Push the commits back to the repo
			// if err = gh.Push(cloneDir); err != nil {
			// 	log.Logger.Errorf("Unable to push change: %s", err)
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// } else {
			// 	log.Logger.Info("Skipping git push due to testing")
			// }
		}
	}
	// Give a generic success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}
