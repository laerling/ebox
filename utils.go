/*
 * Copyright (c) 2018 lærling
 * For details see ./LICENSE
 */

package main

import (
	"errors"
	"os"
	"strings"
)

// sortableStringSlice makes []string sortable. The sorting is case-insensitive.
type sortableStringSlice []string

// implementation of sort.Interface on stringSlice
func (s sortableStringSlice) Len() int {
	return len(s)
}

// implementation of sort.Interface on stringSlice
func (s sortableStringSlice) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

// implementation of sort.Interface on stringSlice
func (s sortableStringSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

// getGithubUser returns the Github usernames for some of the most common Emacs
// distributions.
func getGithubUser(distroName string) (string, error) {

	switch distroName {
	case "doom-emacs":
		return "hlissner", nil
	case "emacs-live":
		return "overtone", nil
	case "prelude":
		return "bbatsov", nil
	case "spacemacs":
		return "syl20bnr", nil
	case "ohai-emacs":
		return "bodil", nil
	}

	return "", errors.New("Could not find Github userName for repository '" + distroName + "'")
}

// makeRepoUrl takes a repository name or github name with repository
// name or URL and returns the fully qualified URL to use with git
// clone, as well as the name of the repository. For example:
// makeRepoUrl("prelude") == ("https://github.com/bbatnov/prelude", "prelude", nil)
// makeRepoUrl("bbatsov/prelude") == ("https://github.com/bbatsov/prelude", "RepoName", nil)
// makeRepoUrl("domain.tld/foo.git") == ("https://domain.tld/foo.git", "foo", nil)
func makeRepoUrl(distroUrlOrRepoName string) (string, string, error) {
	distroUrl := distroUrlOrRepoName

	// if distroUrl does not contain slash, get Github username
	// the exact position of the slash is needed later
	slashIndex := strings.Index(distroUrl, "/")
	if slashIndex < 0 {
		userName, err := getGithubUser(distroUrl)
		if err != nil {
			return "", "", err
		}
		distroUrl = userName + "/" + distroUrl
	}

	// at this point distroUrl can have the form foo/bar, or domain.tld/foo/bar
	// if distroUrl does not contain dot before slash, assume Github
	dotIndex := strings.Index(distroUrl, ".")
	if dotIndex < 0 || dotIndex > slashIndex {
		distroUrl = "github.com/" + distroUrl
	}

	// at this point distroUrl has the form domain.tld/foo/bar
	// extract name of distro
	distroName := distroUrl[strings.LastIndex(distroUrl, "/")+1:]
	if strings.HasSuffix(distroName, ".git") {
		distroName = distroName[:len(distroName)-4]
	}

	return distroUrl, distroName, nil
}

func directoryExists(dirName string) bool {
	dir, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		// fail only if stat failed but file is present
		panic(err)
	}
	// check that it's really a directory, or a symlink pointing to a directory (or to a symlink pointing to... you get the idea)
	if !dir.IsDir() {
		return false
	}
	return true
}

func ensureDirectoryExists(dirName string) error {
	if !directoryExists(dirName) {
		if err := os.Mkdir(dirName, 0755); err != nil {
			return err
		}
	}
	return nil
}

func ensureDirectoryExistsNot(dirName string) error {
	if directoryExists(dirName) {
		if err := os.RemoveAll(dirName); err != nil {
			return err
		}
	}
	return nil
}
