package ebox

import "os"

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
