package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/victims/victims-bot/log"
)

// FileResult which is part of a HashResult
type FileResult struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

// A HashResult from a package submission
type HashResult struct {
	Hash   string       `json:"hash"`
	Name   string       `json:"name"`
	Format string       `json:"format"`
	Files  []FileResult `json:"files"`
}

// SubmitPackage fronts other submittion functions using the proper
// one based on the file suffix
func SubmitPackage(fileName, uri string) (*[]HashResult, error) {
	if strings.HasSuffix(fileName, ".jar") {
		return SubmitJavaPackage(fileName, uri)
	}
	return nil, errors.New("Unknown package type")
}

// SubmitJavaPackage submits a java package for hashing
// Implements https://github.com/victims/victims-java-service
func SubmitJavaPackage(fileName, uri string) (*[]HashResult, error) {
	// NOTE: Hardcoded the OpenShift service url if one is not provided
	if uri == "" {
		uri = "http://victims-java-service.default.svc.cluster.local"
	}
	uri = uri + "/hash"

	result := []HashResult{}
	values := url.Values{}
	values.Set("library2", fileName)

	resp, err := http.PostForm(uri, values)
	if err != nil {
		log.Logger.Errorf("Unable to submit Java Package %s: %s", fileName, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Did not recieve 200. Got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Errorf("Response from server was not able to be read: %s", err)
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Logger.Errorf("Unable to unmarshal server response: %s", err)
		return nil, err
	}
	return &result, nil
}
