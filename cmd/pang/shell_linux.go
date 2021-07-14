package main

import "os/exec"

func openfolder(folder string) error {
	return exec.Command("xdg-open", folder).Start()
}
