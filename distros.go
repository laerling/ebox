package ebox

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func setDistro(homeDirName, distroName string, removeOldSymlink bool) {

	// check that no Emacs instance is running
	if emacsInstanceRunning() {
		fmt.Fprintln(os.Stderr, "Emacs already running. Please terminate it and try again.")
		os.Exit(1)
	}

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

func listDistros(homeDirName, activeDistroPath string) {

	activeDistroDir := path.Base(activeDistroPath)
	// activeDistroDir = ".emacs.d-<distro>"
	// activeDistroDir = "/home/<user>/.emacs.d-<distro>"

	// get active distro if present
	if !strings.HasPrefix(activeDistroDir, distroPrefix) {
		fmt.Fprintln(os.Stderr, "Distribution \""+activeDistroDir+
			"\" must start with \""+distroPrefix+"\". Aborting")
		os.Exit(1)
	}
	activeDistro := activeDistroDir[len(distroPrefix):]

	// open home directory
	homeDir, err := os.Open(homeDirName)
	if err != nil {
		panic(err)
	}

	// read all filenames from home directory
	filenames, err := homeDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	// list distros, highlighting the active one, if present
	activeDistroPresent := false
	for _, filename := range filenames {
		if strings.HasPrefix(filename, ".emacs.d-") {
			filename = filename[len(distroPrefix):]
			if filename == activeDistro {
				fmt.Print("* ")
				activeDistroPresent = true
			} else {
				fmt.Print("  ")
			}
			fmt.Println(filename)
		}
	}

	// warn if active distro not present
	if !activeDistroPresent {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Warning: ~/.emacs.d does not seem to point to a valid distribution")
		fmt.Fprintln(os.Stderr, "         ~/.emacs.d -> "+activeDistroPath)
	}
}
