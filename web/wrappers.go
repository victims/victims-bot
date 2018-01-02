package web

import (
	"io"
	"net/http"
	"reflect"
	"runtime"

	"github.com/victims/victims-bot/log"
)

// getFuncName is a shortcut to get a HandlerFunction's name
func getFuncName(f HandlerFunction) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// RequirePost wrapps a HandlerFunction and only allows POSTs
func RequirePost(f HandlerFunction) HandlerFunction {
	log.Logger.Debugf("Registered %s with mux", getFuncName(f))
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			f(w, req)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Bad Request")
	}
}
