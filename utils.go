package main

import (
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
