package gh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	plumingObject "gopkg.in/src-d/go-git.v4/plumbing/object"

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

	iterations := 0
	defer func() { log.Logger.Debugf("Iterations: %d", iterations) }()

	for {
		iterations = iterations + 1
		file, err := files.Next()
		if err != nil {
			log.Logger.Warn(err)
			return "", err
		}

		if file == nil {
			break
		} else if file.Name == fileName {
			contents, err := file.Contents()
			if err != nil {
				return "", err
			}
			return contents, nil
		}
	}
	return "", errors.New("file not found")
}

// CommitChange commits a change back to the local checkout.
// Returns the hash string of the commit
func CommitChange(repoPath, path string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Logger.Warnf("Unable to open repo at %s. %s", path, err)
		return "", err
	}

	tree, err := repo.Worktree()
	if err != nil {
		log.Logger.Errorf("Unable to open worktree: %s", err)
		return "", err
	}

	hash, err := tree.Add(path)
	if err != nil {
		log.Logger.Errorf("Unable to add %s to worktree: %s", path, err)
		return "", err
	}
	log.Logger.Infof("Locally added %s: %s", path, hash)

	// Info for the commit options
	signature := plumingObject.Signature{
		Name:  "victims-bot",
		Email: "victims-bot@users.noreply.github.com",
		When:  time.Now(),
	}

	// Options passed into tree.Commit
	commitOps := git.CommitOptions{
		Author:    &signature,
		Committer: &signature,
	}
	// The commit message
	commitMsg := fmt.Sprintf("Add hashes to %s", path)
	newHash, err := tree.Commit(commitMsg, &commitOps)
	if err != nil {
		log.Logger.Errorf("Unable to commit %s to worktree: %s", path, err)
		return "", err
	}

	return newHash.String(), nil
}

// Push pushes local changes to the remote repo.
func Push(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Logger.Warnf("Unable to open repo at %s. %s", path, err)
		return err
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
	})

	return err
}
