package web

import (
	"fmt"
	"net/http"

	githubhook "gopkg.in/rjz/githubhook.v0"

	"github.com/victims/victims-bot/cmd"
	"github.com/victims/victims-bot/gh"
	"github.com/victims/victims-bot/log"
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
			_, err = gh.GetContent(cloneDir, commit.ID, file)
			if err != nil {
				log.Logger.Errorf("Error getting contents: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO: Submit to the hash service
			// hashes, err := process.SubmitPackage(fileName, "")
			// if err := process.AddHashesToFile(file, hashes); err != nil {
			// 	log.LoggerInfof("Unable to add hash to file: %s", err)
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			_, err = gh.CommitChange(file)
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
		// Push the commits back to the repo
		if err = gh.Push(cloneDir); err != nil {
			log.Logger.Errorf("Unable to push change: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Give a generic success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}
