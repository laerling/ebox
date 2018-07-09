package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
)

func startDistro(distroDirName, distroName string) {

	// check that the distro exists
	if _, err := os.Stat(distroDirName + "/" + distroName); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "No such distro: '"+distroName+"'")
		os.Exit(1)
	}

	// set this distro as $HOME
	cmd := exec.Command("emacs")
	cmd.Env = append(os.Environ(), "HOME="+distroDirName+"/"+distroName)

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
		panic(err)
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
