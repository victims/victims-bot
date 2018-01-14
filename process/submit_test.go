package process

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	invalidPackage = "test.invalid"
	validPackage   = "test.jar"
)

/*
 * Test helper handlers which respond to the tests as http handlers
 */
// testSuccessHandler fakes a successful hash response
func testSuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `[ {
  "hash" : "3cfc3c06a141ba3a43c6c0a01567dbb7268839f75e6367ae0346bab90507ea09c9ecd829ecec3f030ed727c0beaa09da8c08835a8ddc27054a03f800fa049a0a",
  "name" : "camel-snakeyaml-2.17.4.jar",
  "format" : "SHA512",
  "files" : [ {
    "name" : "org/apache/camel/component/snakeyaml/SnakeYAMLDataFormat",
    "hash" : "cb1e80599bd7de814b63ad699849360b6c5d6dd33b7b7a2da6df753197eee137541c6bfde704c5ab8521e6b7dfb436d57f102f369fc0af36738668e4d1d0ff55"
  } ]
} ]`)
}

// testMalformedHandler returns a malformed response
func testMalformedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `malformed`)
}

// testFailureHandler fakes a server side error
func testFailureHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

// Verify that a []HashResult is returned when a proper response is provided
func TestSubmitPackage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testSuccessHandler))
	defer testServer.Close()

	hashResult, err := SubmitPackage(validPackage, testServer.URL)
	if err != nil {
		t.Errorf("Submission failed: %s\n", err)
		t.Fail()
	}
	if len(*hashResult) != 1 {
		t.Errorf("Expected 1 found %d", len(*hashResult))
		t.Fail()
	}
}

// Verify that if a package type is unknown an error is raised
func TestSubmitPackageWithBadPackage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testSuccessHandler))
	defer testServer.Close()

	invalidPackage := "test.invalid"
	hashResult, err := SubmitPackage(invalidPackage, testServer.URL)
	if err == nil {
		t.Errorf("An error should have been returned for invalid package format %s", invalidPackage)
		t.Fail()
	}
	if hashResult != nil {
		t.Errorf("A result was given when none should have been returned: %#v", hashResult)
		t.Fail()
	}
}

// Verify that if the server responds with a non 200 an error is raised
func TestSubmitPackageWithBadServerResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testFailureHandler))
	defer testServer.Close()

	for _, pkg := range []string{invalidPackage, validPackage} {
		hashResult, err := SubmitPackage(pkg, testServer.URL)
		if err == nil {
			t.Errorf("An error should have been returned for invalid package format %s", invalidPackage)
			t.Fail()
		}
		if hashResult != nil {
			t.Errorf("A result was given when none should have been returned: %#v", hashResult)
			t.Fail()
		}
	}
}

// Verify that if the server responds with malformed data an error is returned
func TestSubmitPackageWithMalformedResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testMalformedHandler))
	defer testServer.Close()

	for _, pkg := range []string{invalidPackage, validPackage} {
		hashResult, err := SubmitPackage(pkg, testServer.URL)
		if err == nil {
			t.Errorf("An error should have been returned for invalid package format %s", invalidPackage)
			t.Fail()
		}
		if hashResult != nil {
			t.Errorf("A result was given when none should have been returned: %#v", hashResult)
			t.Fail()
		}
	}
}
