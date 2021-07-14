package main

import "os/exec"

func openfolder(folder string) error {
	return exec.Command("open", folder).Start()
}
