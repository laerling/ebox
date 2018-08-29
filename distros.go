package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func listDistros(distroDirName string) {

	// open distro directory
	distroDir, err := os.Open(distroDirName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot open directory "+distroDirName)
		fmt.Fprintln(os.Stderr, "Please make sure that it exists and the permissions are set correctly.")
		os.Exit(1)
	}

	// get distros from distro directory
	var distros sortableStringSlice
	distros, err = distroDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	// sort filenames alphabetically (case-insensitively)
	sort.Sort(distros)

	// list distros
	for _, distroName := range distros {
		// check that it's really a distro and doesn't start with a dot
		if directoryExists(distroDirName+PATHSEP+distroName) && !strings.HasPrefix(distroName, ".") {
			fmt.Println(distroName)
		}
	}
}

func downloadOrStartDistro(homeDir, distroDirName, distroName string) {

	// download distro if it does not exist
	distroDir := distroDirName + PATHSEP + distroName
	if !directoryExists(distroDir) {
		if err := downloadDistro(distroDirName, distroName); err != nil {

			// distro not found. Ask user whether to create a new distro
			fmt.Print("Distribution " + distroName + " does not exist. Create it now? (y/N) ")
			var input [1]byte
			_, err = os.Stdin.Read(input[:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading answer")
				os.Exit(1)
			}

			// check answer
			if input[0] == 'y' || input[0] == 'Y' {
				makeNewDistro(homeDir, distroDir, distroName)
			}

			os.Exit(0)
		}

		// if download was successful, return right away
		// instead of starting the distro because the user
		// might want to do some configuration first
		// (e. g. putting proxy variables into init.el) and
		// it's not much overhead to type ebox <distro> again
		// or just press <C-p> on a readline-enabled prompt :P
		return
	}

	// make emacs command
	emacsExe := "emacs"
	if WINDOWS {
		emacsExe += ".exe"
	}

	// set distro as $HOME
	originalDistroDir, err := os.Readlink(distroDir)
	// if readlink fails, assume distroDir is the original dir and continue
	if err != nil {
		originalDistroDir = distroDir
	}
	emacsCmd := exec.Command(emacsExe)
	emacsCmd.Env = append(os.Environ(), "HOME="+originalDistroDir)

	// launch Emacs asynchronously
	err = emacsCmd.Start()
	if err != nil {
		panic(err)
	}
}

func downloadDistro(distroDirName, distroUrlOrName string) error {

	distroUrl, distroName, err := makeRepoUrl(distroUrlOrName)
	if err != nil {
		return err
	}

	// assume https
	if !strings.Contains(distroUrl, "://") {
		distroUrl = "https://" + distroUrl
	}

	// generate git command
	gitExe := "git"
	if WINDOWS {
		gitExe += ".exe"
	}

	// the directory to clone the distro into
	destinationDir := distroDirName + PATHSEP + distroName

	// make git command
	destinationEmacsDir := destinationDir + PATHSEP + ".emacs.d"
	gitCmd := exec.Command(gitExe, "clone", distroUrl, destinationEmacsDir)

	// show git running
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	// run git
	err = gitCmd.Run()
	if err != nil {
		return errors.New("Git error: " + err.Error())
	}

	return nil
}

// makeNewDistro makes the directory for a new distro and fills it with the most
// basic stuff (needed symlinks, environment variable settings for Emacs, ...)
func makeNewDistro(homeDir, distroDir, distroName string) {
	if err := ensureDirectoryExists(distroDir); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot mkdir "+distroName)
		os.Exit(1)
	}

	// symlinks to create
	symlinks := []string{
		".cache",
		".cargo",
		".config",
		".gnupg",
	}

	// create symlinks
	for _, linkName := range symlinks {
		from := homeDir + PATHSEP + linkName
		to := distroDir + PATHSEP + linkName
		err := os.Symlink(from, to)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Warning! Cannot symlink "+from+" to "+to)
		}
	}

	// TODO: Write HTTP(S)_PROXY (if defined) variables to init.el
}
