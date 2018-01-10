package process

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

/*
 * Test helper handlers which respond to the tests as http handlers
 */
// testSuccessJarHandler fakes a successful jar download
func testSuccessJarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/java-archive")
	fmt.Fprintf(w, `data`)
}

// TestGetPackage tests retrieving a package from a remote server
func TestGetPackage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testSuccessJarHandler))
	defer testServer.Close()

	fileName, err := GetPackage(validPackage, testServer.URL)
	defer os.Remove(fileName)

	if err != nil {
		t.Errorf("Error returned getting the package: %s", err)
		t.Fail()
	}

	if fileName == "" {
		t.Errorf("No file name was returned")
		t.Fail()
	}
}

// TestGetPackageWithServerError ensures an error returns when the server fails
func TestGetPackageWithServerError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testFailureHandler))
	defer testServer.Close()

	fileName, err := GetPackage(validPackage, testServer.URL)
	if err == nil {
		t.Errorf("Error should have been returned")
		t.Fail()
	}

	if fileName != "" {
		t.Errorf("Filename was returned but should not have been: %s", fileName)
		t.Fail()
	}
}

// TestGetPackageWithBadURL ensures an error returns when the URL is not accepted
func TestGetPackageWithBadURL(t *testing.T) {
	fileName, err := GetPackage(validPackage, "baduritest://localhost")
	if err == nil {
		t.Errorf("Error should have been returned")
		t.Fail()
	}
	if fileName != "" {
		t.Errorf("Filename was returned but should not have been: %s", fileName)
		t.Fail()
	}
}
