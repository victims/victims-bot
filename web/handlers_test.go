package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/victims/victims-bot/cmd"
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

// TestHookWithImproperBadData verifies Hook errors if improper data is passed
func TestHookWithImproperBadData(t *testing.T) {
	// Manually set the secret
	cmd.Config.Secret = "test"
	// Create the test server
	testServer := httptest.NewServer(http.HandlerFunc(Hook))
	defer testServer.Close()

	// Post without any data
	resp, err := http.Post(testServer.URL, "application/json", nil)
	if err != nil {
		t.Errorf("Errored trying to create ping: %s", err)
		t.Fail()
	}

	// The result should be ISE
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected %d got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

// TestHookWithGoodPing verifies good pings make it through and repsond properly
func TestHookWithGoodPing(t *testing.T) {
	// Manually set the secret
	cmd.Config.Secret = "test"
	// Create the test server
	testServer := httptest.NewServer(http.HandlerFunc(Hook))
	defer testServer.Close()

	// Set up the request
	postData, _ := ioutil.ReadFile("testdata/ping.json")
	req, err := http.NewRequest("POST", testServer.URL, bytes.NewReader(postData))
	if err != nil {
		t.Errorf("Errored trying to create ping: %s", err)
		t.Fail()
	}
	req.Header.Set("X-Hub-Signature", "sha1=5e13e22615a472ab78d659ed9482478faab4a25d")
	req.Header.Set("X-GitHub-Event", "ping")
	req.Header.Set("X-Github-Delivery", "72d3162e-cc78-11e3-81ab-4c9367dc0958")

	// Submit the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Errored trying to post ping: %s", err)
		t.Fail()
	}

	// The result should be OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestHookWithPushEvent verifies the process when a good hook is sent
func TestHookWithPushEvent(t *testing.T) {
	os.Setenv("VICTIMS_BOT_TEST", "1")
	cmd.Config.GitRepo = "https://github.com/victims/victims-cve-db.git"
	// Manually set the secret
	cmd.Config.Secret = "test"
	// Create the test server
	testServer := httptest.NewServer(http.HandlerFunc(Hook))
	defer testServer.Close()

	// Set up the request
	postData, _ := ioutil.ReadFile("testdata/push.json")
	req, err := http.NewRequest("POST", testServer.URL, bytes.NewReader(postData))
	if err != nil {
		t.Errorf("Errored trying to create push: %s", err)
		t.Fail()
	}
	req.Header.Set("X-Hub-Signature", "sha1=7c487f2cf15fa372b29621ad5e5d57d52c9d98d4")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Github-Delivery", "72d3162e-cc78-11e3-81ab-4c9367dc0958")
	// Submit the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Errored trying to post ping: %s", err)
		t.Fail()
	}

	// The result should be OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusOK, resp.StatusCode)
	}
}
