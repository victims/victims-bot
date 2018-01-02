package gh

import (
	"errors"
	"io/ioutil"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

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

// GetContent returns the content of a specific file at a specific commit
func GetContent(path, hash, fileName string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Logger.Warnf("Unable to open repo at %s. %s", path, err)
		return "", err
	}

	commit, err := repo.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		log.Logger.Warnf("Unable to find commit %s in %s. %s", hash, path, err)
		return "", err
	}

	files, err := commit.Files()
	if err != nil {
		log.Logger.Warnf("Unable to list files for %s in %s. %s", hash, path, err)
		return "", err
	}

	for {
		file, err := files.Next()
		if err != nil {
			log.Logger.Warn(err)
			return "", err
		}

		if file == nil {
			break
		} else if file.Name == fileName {
			contents, _ := file.Contents()
			return contents, nil
		}
	}
	return "", errors.New("file not found")
}
