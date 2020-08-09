package main

import "os/exec"

func openfolder(folder string) {
	exec.Command("open", folder).Start()
}
