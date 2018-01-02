package web

import "net/http"

// HandlerFunction represents the function type for handler functions
type HandlerFunction func(http.ResponseWriter, *http.Request)
