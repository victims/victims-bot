package main

import (
	"net/http"

	"github.com/victims/victims-bot/cmd"
	"github.com/victims/victims-bot/log"
	"github.com/victims/victims-bot/web"
)

// main is the main entry point
func main() {
	cmd.ParseFlags()

	mux := http.NewServeMux()
	log.Logger.Debug("Mux created")
	mux.HandleFunc("/", web.RequirePost(web.Hook))

	s := &http.Server{
		Addr:    cmd.Config.BindTo,
		Handler: mux,
	}
	log.Logger.Infof("Serving on %s\n", cmd.Config.BindTo)
	log.Logger.Fatal(s.ListenAndServe())
}
