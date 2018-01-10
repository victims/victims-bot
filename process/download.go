package process

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/victims/victims-bot/log"
)

// GetPackage downloads a package from a remote location
func GetPackage(name, url string) (string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "http://") {
		return "", errors.New("URL must be http:// or https://")
	}
	file, err := ioutil.TempFile("", name)
	defer file.Close()
	if err != nil {
		log.Logger.Errorf("Unable create temp file %s", err)
		return "", err
	}

	timeout := time.Duration(30 * time.Second)
	log.Logger.Debugf("Set HTTP client timeout to %f seconds", timeout.Seconds())
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Logger.Errorf("Unable to download %s: %s", url, err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Did not recieve 200. Got %d", resp.StatusCode)
	}

	copied, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Logger.Errorf("Unable to copy data to file %s: %s", file.Name(), err)
		return "", err
	}
	log.Logger.Infof("Copied %d to %s", copied, file.Name())
	return file.Name(), nil
}
