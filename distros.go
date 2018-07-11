package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func downloadOrStartDistro(distroDirName, distroName string) {

	// download distro if it does not exist
	if _, err := os.Stat(distroDirName + PATHSEP + distroName); os.IsNotExist(err) {
		if err = downloadDistro(distroDirName, distroName); err != nil {
			fmt.Fprintln(os.Stderr, "Error downloading distro '"+distroName+"':")
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		// if download was successful return right away
		// instead of starting the distro because the user
		// might want to do some configuration first
		// (e. g. putting proxy variables into init.el) and
		// it's not much overhead to type ebox <distro> again
		// or just press <C-p> on a readline-enabled prompt :P
		return
	}

	// set this distro as $HOME
	emacsExe := "emacs"
	if WINDOWS {
		emacsExe += ".exe"
	}

	// make emacs command
	emacsCmd := exec.Command(emacsExe)
	emacsCmd.Env = append(os.Environ(), "HOME="+distroDirName+PATHSEP+distroName)

	// launch Emacs asynchronously
	err := emacsCmd.Start()
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

func downloadDistro(distroDirName, distroUrlOrName string) error {

	distroUrl, distroName, err := makeRepoUrl(distroUrlOrName)
	if err != nil {
		return err
	}

	// assume https
	if !strings.Contains(distroUrl, "://") {
		distroUrl = "https://" + distroUrl
	}

	// generate git command
	gitExe := "git"
	if WINDOWS {
		gitExe += ".exe"
	}

	// the directory to clone the distro into
	destinationDir := distroDirName + PATHSEP + distroName

	// make git command
	destinationEmacsDir := destinationDir + PATHSEP + ".emacs.d"
	gitCmd := exec.Command(gitExe, "clone", distroUrl, destinationEmacsDir)

	// show git running
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	// run git
	err = gitCmd.Run()
	if err != nil {
		return errors.New("Git error: " + err.Error())
	}

	return nil
}
