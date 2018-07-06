package main

import (
	"fmt"
	"os"
	"os/user"
)

const distroPrefix = ".emacs.d-"

func main() {

	// get ~/.emacs.d
	usr, _ := user.Current()
	homeDirName := usr.HomeDir
	emacsLinkName := homeDirName + "/.emacs.d"

	// check existence of ~/.emacs.d
	// Just because Readlink succeeds doesn't mean that the links' destination exists!
	emacsDir, readLinkErr := os.Readlink(emacsLinkName)
	// But it does mean that the link itself exists.
	emacsLinkExists := readLinkErr == nil

	// exit if non-empty argument supplied and ~/.emacs.d exists, but is not a symlink
	// non-empty argument means that we'll switch to the distro with that argument as name
	if len(os.Args) > 1 && len(os.Args[1]) > 0 && emacsLinkExists && readLinkErr != nil {
		fmt.Fprintln(os.Stderr, emacsLinkName, "is not a symbolic link. Aborting")
		os.Exit(1)
	}

	// if argument supplied, set distro, else list existing distros
	if len(os.Args) > 1 {
		setDistro(homeDirName, os.Args[1], emacsLinkExists)
	} else {
		listDistros(homeDirName, emacsDir)
	}
}
