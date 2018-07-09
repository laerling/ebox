package main

import (
	"errors"
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
	case "spacemacs":
		return "syl20bnr", nil
	case "prelude":
		return "bbatsov", nil
	case "doom-emacs":
		return "hlissner", nil
		/*
			case "emacs-live":
				return "overtone", nil
		*/
	}

	return "", errors.New("Could not find userName for " + distroName)
}
