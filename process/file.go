package process

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/victims/victims-bot/log"
)

// Affected represents an affected entity in CVEDBEntry
type Affected struct {
	GroupID    string   `yaml:"groupId,omitempty"`
	ArtifactID string   `yaml:"artifactId,omitempty"`
	Version    []string `yaml:"version,omitempty"`
	FixedIn    []string `yaml:"fixedin,omitempty"`
}

// CVEDBEntry represents a single entry in the cvedb
type CVEDBEntry struct {
	CVE         string       `yaml:"cve"`
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	CVSSV2      string       `yaml:"cvss_v2"`
	References  []string     `yaml:"references"`
	Hash        string       `yaml:"hash"`
	FileHashes  []FileResult `yaml:"file_hashes"`
	Affected    []Affected   `yaml:"affected"`
	URL         string       `yaml:"url"`
	Name        string       `yaml:"name"`
}

// AddHashesToFile adds hash information to a cvedb file
func AddHashesToFile(fileName string, hashes HashResult) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Logger.Errorf("Unable to read file %s: %s", fileName, err)
		return err
	}

	// Unmarshal file into a struct
	cvedbEntry := CVEDBEntry{}
	if err = yaml.Unmarshal(data, &cvedbEntry); err != nil {
		log.Logger.Errorf("Unable to deserialize %s: %s", fileName, err)
		return err
	}

	// Add the hash info
	cvedbEntry.Hash = hashes.Hash
	cvedbEntry.FileHashes = hashes.Files

	// Marshal the contents back and write the results to disk
	result, err := yaml.Marshal(cvedbEntry)
	if err != nil {
		log.Logger.Errorf("Unable to serialize %s: %s", fileName, err)
		return err
	}

	// Write the file back
	if err = ioutil.WriteFile(fileName, result, os.ModeExclusive); err != nil {
		log.Logger.Errorf("Unable to write update to file %s: %s", fileName, err)
		return err
	}

	return nil
}
