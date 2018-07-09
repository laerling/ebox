package main

import (
	"os"
	"os/user"
)

func main() {

	// get distro directory name
	usr, _ := user.Current()
	distroDirName := usr.HomeDir + PATHSEP + "emacs"

	// if distro directory does not exist, create it
	if _, err := os.Stat(distroDirName); os.IsNotExist(err) {
		os.Mkdir(distroDirName, 0755)
	}

	// if argument supplied, start distribution, else list existing distributions
	if len(os.Args) > 1 {
		startDistro(distroDirName, os.Args[1])
	} else {
		listDistros(distroDirName)
	}
}
