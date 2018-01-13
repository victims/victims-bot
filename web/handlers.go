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
	hook, err := githubhook.Parse(cmd.Config.GetSecret(), req)
	if err != nil {
		log.Logger.Infof("could not read webhook: %s\n", err)
		return
	}
	//body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.Errorf("Web hook unable to read the request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pushEvent := gh.PushEvent{}
	//json.Unmarshal(body, &pushEvent)
	hook.Extract(&pushEvent)
	url := "https://github.com/victims/victims-cve-db.git" //"git@github.com:victims/victims-cve-db.git"
	cloneDir, err := gh.Clone(url)
	if err != nil {
		log.Logger.Errorf("Web hook failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, cloneDir)

	// NOTE: that if anything fails while processing the change the entire
	// set of work will not be pushed back to the repo
	for _, commit := range pushEvent.Commits {
		// Probably put this in it's own bounded goroutine
		for _, file := range commit.Added {
			_, err := gh.GetContent(cloneDir, commit.ID, file)
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
		}
	}

	if err = gh.Push(cloneDir); err != nil {
		log.Logger.Errorf("Unable to push change: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")

}
