package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

func setDistro(homeDirName, distroName string, removeOldSymlink bool) {

	// check that no Emacs instance is using the current distro
	blockingPid, err := checkNoEmacsBlocking()
	if err != nil {
		fmt.Fprintln(os.Stderr, "An Emacs instance with the PID "+
			blockingPid+" is using the active distribution.")
		fmt.Fprintln(os.Stderr, "Please terminate it and try again.")
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

	// filter distributions
	distros := make(sortableStringSlice, 0, len(filenames))
	for _, filename := range filenames {
		if strings.HasPrefix(filename, ".emacs.d-") {
			distros = append(distros, filename)
		}
	}

	// sort filenames alphabetically (case-insensitively)
	sort.Sort(distros)

	// list distros, highlighting the active one, if present
	activeDistroPresent := false
	for _, distro := range distros {
		distro = distro[len(distroPrefix):]
		if distro == activeDistro {
			fmt.Print("* ")
			activeDistroPresent = true
		} else {
			fmt.Print("  ")
		}
		fmt.Println(distro)
	}

	// warn if active distro not present
	if !activeDistroPresent {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Warning: ~/.emacs.d does not seem to point to a valid distribution")
		fmt.Fprintln(os.Stderr, "         ~/.emacs.d -> "+activeDistroPath)
	}
}
