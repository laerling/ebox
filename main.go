package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {

	// get distro directory name
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := usr.HomeDir
	distroDirName := homeDir + PATHSEP + "emacs"

	// if distro directory does not exist, create it
	if err := ensureDirectoryExists(distroDirName); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot mkdir "+distroDirName)
		os.Exit(1)
	}

	// if argument supplied, start distribution, else list existing distributions
	if len(os.Args) > 1 {
		downloadOrStartDistro(homeDir, distroDirName, os.Args[1])
	} else {
		listDistros(distroDirName)
	}
}
