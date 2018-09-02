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

	// list distros, if no argument supplied
	if len(os.Args) <= 1 {
		panicOnError(listDistros(distroDirName))
		os.Exit(0)
	}

	// start distro, if its directory exists
	distroDir := distroDirName + PATHSEP + os.Args[1]
	if directoryExists(distroDir) {
		startDistro(homeDir, distroDir)
		os.Exit(0)
	}

	// create the distribution, if its directory does not already exist
	panicOnError(createDistro(homeDir, distroDirName, os.Args[1]))
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
