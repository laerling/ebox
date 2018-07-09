package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func downloadOrStartDistro(distroDirName, distroName string) {

	// download distro if it does not exist
	if _, err := os.Stat(distroDirName + PATHSEP + distroName); os.IsNotExist(err) {
		if err := downloadDistro(distroDirName, distroName); err != nil {
			fmt.Fprintln(os.Stderr, "No such distro: '"+distroName+"' "+
				"and no downloadable distro found in list.")
			os.Exit(1)
		}

		// don't start the distro because the user might want to do some
		// configuration first and it's not much overhead to type ebox <distro>
		// again or just press <C-p> on a readline-enabled prompt :P

		return
	}

	// set this distro as $HOME
	emacsExecutable := "emacs"
	if WINDOWS {
		emacsExecutable += ".exe"
	}

	// make emacs command
	cmd := exec.Command(emacsExecutable)
	cmd.Env = append(os.Environ(), "HOME="+distroDirName+PATHSEP+distroName)

	// launch Emacs asynchronously
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

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
	for _, distro := range distros {
		fmt.Println(distro)
	}
}

func downloadDistro(distroDirName, distroName string) error {

	// if distroName does not contain slash, get Github username
	slashIndex := strings.Index(distroName, "/")
	if slashIndex < 0 {
		userName, err := getGithubUser(distroName)
		if err != nil {
			return err
		}
		distroName = userName + "/" + distroName
	}

	// distroName is guaranteed to have the form (domain.tld/)?foo/bar now

	// if distroName does not contain dot before slash, assume github
	dotIndex := strings.Index(distroName, ".")
	if dotIndex < 0 || dotIndex > slashIndex {
		distroName = "github.com/" + distroName
	}

	// distroName is guaranteed to have the form domain.tld/foo/bar now

	// generate URL to clone and name of distro
	distroName = "https://" + distroName
	distroNameBase := distroName[strings.LastIndex(distroName, "/")+1:]

	// generate git command
	gitExecutable := "git"
	if WINDOWS {
		gitExecutable += ".exe"
	}

	// make sure the destination directory exists
	destinationDir := distroDirName + PATHSEP + distroNameBase
	err := os.Mkdir(destinationDir, 0755)
	if err != nil {
		// TODO don't fail if dir already exists
		fmt.Fprintln(os.Stderr, "Cannot mkdir "+destinationDir)
		os.Exit(1)
	}

	// run git
	destinationEmacsDir := destinationDir + PATHSEP + ".emacs.d"
	fmt.Println("Running " + gitExecutable + " clone " + distroName + " " + destinationEmacsDir)
	cmd := exec.Command(gitExecutable, "clone", distroName, destinationEmacsDir)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Git error: "+err.Error())
		return err
	}

	return nil
}
