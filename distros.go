package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func listDistros(distroDirName string) error {

	// open distro directory
	distroDir, err := os.Open(distroDirName)
	if err != nil {
		return errors.New("Cannot open directory " +
			distroDirName +
			". Please make sure that it exists and the permissions are set correctly.")
	}

	// get distros from distro directory
	var distros sortableStringSlice
	distros, err = distroDir.Readdirnames(0)
	if err != nil {
		return err
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

	return nil
}

func startDistro(homeDir, distroDir string) {

	// set distro as $HOME
	originalDistroDir, err := os.Readlink(distroDir)
	// if readlink fails, assume distroDir is the original dir and continue
	if err != nil {
		originalDistroDir = distroDir
	}

	// check symlinks are present in distro
	ensureSymlinksPresent(homeDir, distroDir)

	// make Emacs command
	emacsExe := "emacs"
	if WINDOWS {
		emacsExe += ".exe"
	}
	emacsCmd := exec.Command(emacsExe)
	emacsCmd.Env = append(os.Environ(), "HOME="+originalDistroDir)

	// launch Emacs asynchronously
	err = emacsCmd.Start()
	if err != nil {
		panic(err)
	}
}

func createDistro(homeDir, distroDirName, distroName string) error {
	distroDir := distroDirName + PATHSEP + distroName
	if err := downloadDistro(distroDirName, distroName); err == nil {
		return nil
	}

	// distro not found. Ask user whether to create a new distro
	fmt.Print("Distribution " + distroName + " does not exist. Create it now? (y/N) ")
	var input [1]byte
	_, err := os.Stdin.Read(input[:])
	if err != nil {
		return errors.New("Error reading answer")
	}

	// check answer
	if input[0] != 'y' && input[0] != 'Y' {
		return nil
	}

	// create new distro
	if err := ensureDirectoryExists(distroDir); err != nil {
		return errors.New("Cannot mkdir " + distroName)
	}

	// prepare distro
	makeInitFile(distroDir)

	return nil
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

func ensureSymlinksPresent(homeDir, distroDir string) {

	// symlinks to create
	symlinks := []string{
		".cache",
		".cargo",
		".config",
		".gitconfig",
		".gnupg",
		".local",
		".mozilla",
		".rustup",
	}

	// create symlinks
	for _, linkName := range symlinks {

		source := homeDir + PATHSEP + linkName
		target := distroDir + PATHSEP + linkName

		// skip if source does not exist
		_, err := os.Stat(source)
		if err != nil {
			continue
		}

		// make symlink
		if err := os.Symlink(source, target); err == nil {
			fmt.Println("Created symlink " + target + " -> " + source)
		}
	}
}

func makeInitFile(distroDir string) {

	// get environment variables
	envHttpProxy := os.Getenv("http_proxy")
	envHttpsProxy := os.Getenv("https_proxy")

	// make string to write to init.el
	initElString := ""

	// append http_proxy
	if envHttpProxy != "" {
		envHttpProxySplitted := strings.Split(envHttpProxy, "://")
		initElString += "\n     (\"" +
			envHttpProxySplitted[0] + "\" . \"" +
			envHttpProxySplitted[1] + "\")"
	}

	// append https_proxy
	if envHttpsProxy != "" {
		envHttpsProxySplitted := strings.Split(envHttpsProxy, "://")
		initElString += "\n     (\"" +
			envHttpsProxySplitted[0] + "\" . \"" +
			envHttpsProxySplitted[1] + "\")"
	}

	// return if nothing to write
	if initElString == "" {
		return
	}

	// make .emacs.d directory
	emacsDir := distroDir + PATHSEP + ".emacs.d"
	ensureDirectoryExists(emacsDir)

	// write init.el
	stringToWrite := []byte("(set 'url-proxy-services '(" + initElString + "))\n\n")
	initFileName := emacsDir + PATHSEP + "init.el"
	if err := ioutil.WriteFile(initFileName, stringToWrite, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot write to init.el."+
			" This could make the distribution unusable."+
			" Please check the init.el file manually")
	}
}
