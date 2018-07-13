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
const githubUser string = "bodil"
const githubUserN string = "nonexistent"

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

	os.Exit(m.Run())
}

// runs main and checks if distroDir exists
func runMain(t *testing.T, expectSuccess, expectDistroDir bool, distroName, arg string) {

	fmt.Println("Running with argument '" + arg + "'")

	// make sure distro does not exist
	t.Log("Making sure distribution directory does not exist: '" + distroDir + "'")
	if _, err := os.Stat(distroDir); err == nil {
		if err = os.RemoveAll(distroDir); err != nil {
			t.Fatal("Cannot remove distribution " + distroName + ": " + err.Error())
		}
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

// foo
func TestDownloadName(t *testing.T) {
	runMain(t, true, true, distroName, distroName)
}

// foo/bar
func TestDownloadSlash(t *testing.T) {
	runMain(t, true, true, distroName, githubUser+"/"+distroName)
}

// domain.tld/foo
func TestDownloadDomain(t *testing.T) {
	runMain(t, true, true, distroName, "github.com/"+githubUser+"/"+distroName)
}

// nonexistent
func TestDownloadNameN(t *testing.T) {
	runMain(t, false, false, distroNameN, distroNameN)
}

// nonexistent/nonexistent
func TestDownloadSlashN(t *testing.T) {
	runMain(t, false, false, distroNameN, githubUserN+"/"+distroNameN)
}

// domain.tld/nonexistent/nonexistent
func TestDownloadDomainN(t *testing.T) {
	runMain(t, false, false, distroNameN, "domain.tld/"+githubUserN+"/"+distroNameN)
}
