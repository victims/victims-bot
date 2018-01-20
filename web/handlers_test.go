package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/victims/victims-bot/gh"

	githubhook "gopkg.in/rjz/githubhook.v0"
)

// TestPingEvent verifies pings respond properly
func TestPingEvent(t *testing.T) {
	event := gh.PingEvent{
		Hook:   "ping",
		HookId: "ping",
		Zen:    "This may all be a simulation",
	}
	payload, _ := json.Marshal(event)
	hook := githubhook.Hook{
		Event:     "ping",
		Id:        "ping",
		Payload:   payload,
		Signature: "",
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/ping", nil)
	// Execue the event
	pingEvent(&hook, recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusOK, recorder.Code)
	}
}

// TestPingEventWithBadData verifies bad data returns ISE
func TestPingEventWithBadData(t *testing.T) {
	hook := githubhook.Hook{
		Event:     "ping",
		Id:        "ping",
		Payload:   []byte("BAAAAD"),
		Signature: "",
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/ping", nil)
	// Execue the event
	pingEvent(&hook, recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected %d got %d", http.StatusInternalServerError, recorder.Code)
	}
}
