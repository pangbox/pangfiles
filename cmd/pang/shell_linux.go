package main

import "os/exec"

func openfolder(folder string) {
	exec.Command("xdg-open", folder).Start()
}
