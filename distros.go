package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
)

func startDistro(distroDirName, distroName string) {

	// check that the distro exists
	if _, err := os.Stat(distroDirName + PATHSEP + distroName); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "No such distro: '"+distroName+"'")
		os.Exit(1)
	}

	// set this distro as $HOME
	executable := "emacs"
	if WINDOWS {
		executable += ".exe"
	}

	// make emacs command
	cmd := exec.Command(executable)
	cmd.Env = append(os.Environ(), "HOME="+distroDirName+PATHSEP+distroName)

	// launch Emacs asynchronously
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func listDistros(distroDirName string) {

	// open distro directory
	distroDir, err := os.Open(distroDirName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot open directory "+distroDirName)
		fmt.Fprintln(os.Stderr, "Please make sure that it exists and the permissions are set correctly.")
		os.Exit(1)
	}

	// get distros from distro directory
	var distros sortableStringSlice
	distros, err = distroDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	// sort filenames alphabetically (case-insensitively)
	sort.Sort(distros)

	// list distros
	for _, distro := range distros {
		fmt.Println(distro)
	}
}
