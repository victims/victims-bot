package gh

import (
	"io/ioutil"

	"gopkg.in/src-d/go-git.v4"

	"github.com/Sirupsen/logrus"
	"github.com/victims/victims-bot/log"
)

// Clone clones a remote git repository
func Clone(url string) (string, error) {
	cloneDir, err := ioutil.TempDir("", "clone")
	if err != nil {
		log.Logger.Warnf("Unable to create temp dir %s. %s", cloneDir, err)
		return "", err
	}
	log.Logger.Debugf("Cloning to %s\n", cloneDir)

	repo, err := git.PlainClone(
		cloneDir, false, &git.CloneOptions{
			URL:               url,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

	if err != nil {
		log.Logger.Warnf("Unable to clone to %s. %s", cloneDir, err)
		return "", err
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := repo.Head()

	if err != nil {
		log.Logger.Warnf("Unable to retrieve HEAD. %s", err)
		return "", err
	}

	// ... retrieving the commit object
	commit, err := repo.CommitObject(ref.Hash())

	if err != nil {
		log.Logger.Warnf("Unable to retrieve the latest commit. %s", err)
		return "", err
	}

	log.Logger.Infof("Cloned %s to %s", url, cloneDir)
	log.Logger.WithFields(
		logrus.Fields{
			"hash":    commit.Hash.String(),
			"author":  commit.Author.Name,
			"email":   commit.Author.Email,
			"message": commit.Message,
		}).Info()
	return cloneDir, nil
}
