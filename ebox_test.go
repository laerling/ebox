package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"testing"
)

const programName string = "ebox"

// N as a suffix for a variable or a constant stands for Negative

const distroName string = "ohai-emacs"
const distroNameN string = "nonexistent"
const remoteUser string = "bodil"
const remoteUserN string = "nonexistent"

var distroDir string
var distroDirN string

func TestMain(m *testing.M) {

	// set global variables
	usr, err := user.Current()
	if err != nil {
		panic("Error getting current user: " + err.Error())
	}
	distroDir = usr.HomeDir + PATHSEP + "emacs" + PATHSEP + distroName
	distroDirN = usr.HomeDir + PATHSEP + "emacs" + PATHSEP + distroNameN

	// run tests
	exitState := m.Run()

	// clean up
	if err := ensureDirectoryNotExists(distroDirN); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot delete directory "+distroDirN)
		os.Exit(1)
	}

	// exit
	os.Exit(exitState)
}

// runs main and checks if distroDir exists
func runMain(t *testing.T, expectSuccess, expectDistroDir bool, distroName, arg string) {

	fmt.Println("Running with argument '" + arg + "'")

	// make sure distro for positive tests does not exist
	t.Log("Making sure distribution directory does not exist: '" + distroDir + "'")
	if err := ensureDirectoryNotExists(distroDir); err != nil {
		t.Fatal("Cannot remove distribution " + distroName + ": " + err.Error())
	}

	// make sure distro for negative tests does not exist
	t.Log("Making sure distribution directory does not exist: '" + distroDirN + "'")
	if err := ensureDirectoryNotExists(distroDirN); err != nil {
		t.Fatal("Cannot remove distribution " + distroNameN + ": " + err.Error())
	}

	// set command line argument
	os.Args = []string{"", arg}
	t.Log("Program arguments: '" + strings.Join(os.Args, "', '") + "'")
	t.Log("Expect distribution directory to exist after run: " + strconv.FormatBool(expectDistroDir))

	// run program
	cmd := exec.Command(programName, arg)
	err := cmd.Run()
	t.Log("Program finished")
	if cmd.ProcessState.Success() != expectSuccess {
		t.Fatal("Wrong exit status: ", err)
	}

	// check that distro dir exists
	_, err = os.Stat(distroDir)
	distroDirExists := !os.IsNotExist(err)
	if distroDirExists != expectDistroDir {
		t.Error("Distribution directory exists, expected " +
			strconv.FormatBool(expectDistroDir) + ", found " +
			strconv.FormatBool(distroDirExists) + ".")
	}
}

/*
 * negative tests
 */

// nonexisting repo of the form 'nonexistent'
func TestDownloadNameN(t *testing.T) {
	runMain(t, false, false, distroNameN, distroNameN)
}

// there is no test for the remote 'nonexisting/nonexisting', because git would just assume
// the repo to be a private one on Github, thus asking for a username and password

// nonexisting remote of the form 'domain.tld/foo'
func TestDownloadDomainN(t *testing.T) {
	runMain(t, false, false, distroNameN, "domain.tld/"+remoteUserN+"/"+distroNameN)
}

/*
 * positive tests
 * the positive tests run last, because I want to keep the distro
 */

// existing repo of the form foo
func TestDownloadName(t *testing.T) {
	runMain(t, true, true, distroName, distroName)
}

// existing remote of the form foo/bar
func TestDownloadSlash(t *testing.T) {
	runMain(t, true, true, distroName, remoteUser+"/"+distroName)
}

// existing remote of the form 'domain.tld/foo'
func TestDownloadDomain(t *testing.T) {
	runMain(t, true, true, distroName, "github.com/"+remoteUser+"/"+distroName)
}
