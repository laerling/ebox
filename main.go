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

	// if argument supplied, start or make distribution
	if len(os.Args) > 1 {
		distroDir := distroDirName + PATHSEP + os.Args[1]
		if directoryExists(distroDir) {
			startDistro(homeDir, distroDir)
		} else {
			makeDistro(homeDir, distroDirName, os.Args[1])
		}
	} else {
		// else list existing distributions
		listDistros(distroDirName)
	}
}
