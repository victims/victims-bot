package web

import (
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

	pushEvent := gh.PushEvent{}
	hook.Extract(&pushEvent)
	//TODO
}
