package main

import (
	"os"
	"os/user"
)

func main() {

	// get distribution directory name
	usr, _ := user.Current()
	distroDirName := usr.HomeDir + PATHSEP + "emacs"

	// if argument supplied, start distribution, else list existing distributions
	if len(os.Args) > 1 {
		startDistro(distroDirName, os.Args[1])
	} else {
		listDistros(distroDirName)
	}
}
