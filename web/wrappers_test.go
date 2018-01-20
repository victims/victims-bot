package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler(w http.ResponseWriter, req *http.Request) {}

// TestRequirePostNoPost verifies wrapping functionality of handlers rejects non POST
func TestRequirePostWithGet(t *testing.T) {
	handler := RequirePost(testHandler)
	for _, method := range []string{"GET", "PUT", "DELETE", "HEAD", "TRACE"} {
		request := httptest.NewRequest(method, "/", nil)
		recorder := httptest.NewRecorder()
		handler(recorder, request)
		if recorder.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected %d got %d", http.StatusMethodNotAllowed, recorder.Code)
		}
	}
}

// TestRequirePostWithPost verifies wrapping functionality of handlers allows POST
func TestRequirePostWithPost(t *testing.T) {
	handler := RequirePost(testHandler)
	request := httptest.NewRequest("POST", "/", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusOK, recorder.Code)
	}
}
