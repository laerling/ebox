package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
)

const distroPrefix = ".emacs.d-"

func main() {

	// get ~/.emacs.d
	usr, _ := user.Current()
	homeDirName := usr.HomeDir
	emacsLinkName := homeDirName + "/.emacs.d"

	// check ~/.emacs.d
	emacsLinkStat, _ := os.Stat(emacsLinkName)
	emacsLinkExists := emacsLinkStat != nil
	// if ~/.emacs.d is a symlink, it is not guaranteed that its destination exists!
	emacsDir, readLinkErr := os.Readlink(emacsLinkName)

	// exit if non-empty argument supplied and ~/.emacs.d exists, but is not a symlink
	// non-empty argument means that we'll switch to the distro with that argument as name
	if len(os.Args) > 1 && len(os.Args[1]) > 0 && emacsLinkExists && readLinkErr != nil {
		fmt.Fprintln(os.Stderr, emacsLinkName, "is not a symbolic link. Aborting")
		os.Exit(1)
	}

	// list distros if no arguments supplied
	if len(os.Args) == 1 {
		listDistros(homeDirName, path.Base(emacsDir))
		return
	}

	// else, set a distro
	setDistro(homeDirName, os.Args[1], emacsLinkExists)
}

func listDistros(homeDirName, activeDistroDir string) {

	var currentDistro string

	// get current distro if present
	if activeDistroDir != "" /* link has been read */ {
		if !strings.HasPrefix(activeDistroDir, distroPrefix) {
			fmt.Fprintln(os.Stderr, "Distribution \""+activeDistroDir+
				"\" must start with \""+distroPrefix+"\". Aborting")
			os.Exit(1)
		}
		currentDistro = activeDistroDir[len(distroPrefix):]
	}

	// list distros, highlighting the current one, if present
	homeDir, err := os.Open(homeDirName)
	if err != nil {
		panic(err)
	}
	filenames, err := homeDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	for _, filename := range filenames {
		if strings.HasPrefix(filename, ".emacs.d-") {
			filename = filename[len(distroPrefix):]
			if filename == currentDistro {
				fmt.Print("* ")
			} else {
				fmt.Print("  ")
			}
			fmt.Println(filename)
		}
	}
}

func setDistro(homeDirName, distroName string, removeOldSymlink bool) {

	// check that distro exists
	distroDirName := homeDirName + "/.emacs.d-" + distroName
	if _, err := os.Stat(distroDirName); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "No such distro: '"+distroName+"'")
		os.Exit(1)
	}

	// Remove old symlink
	emacsLink := homeDirName + "/.emacs.d"
	if removeOldSymlink {
		err := os.Remove(emacsLink)
		if err != nil {
			panic(err)
		}
	}

	// Create new symlink
	os.Symlink(distroDirName, emacsLink)
}
