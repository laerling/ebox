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

func emacsInstanceRunning() bool {
	/*
		Though we could simply run `ps` and search for "emacs" in the output,
		this program strives to be as dependency-free as possible.
		This shouldn't be (much) less performant than ps anyway.
	*/

	// open /proc
	procDir, err := os.Open("/proc")
	if err != nil {
		panic(err)
	}

	// read all filenames from home directory
	filenames, err := procDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	// read every /proc/*/comm file. That's what ps does. That's what heroes do.
	for _, filename := range filenames {

		// try opening /proc/<filename>/comm
		commFile, _ := os.Open("/proc/" + filename + "/comm")
		if commFile == nil {
			continue
		}

		// read process name from comm file
		emacsProcessName := "emacs"
		commBuffer := make([]byte, len(emacsProcessName))
		commFile.Read(commBuffer)

		// check name
		if string(commBuffer) == emacsProcessName {
			return true
		}
	}

	return false
}
