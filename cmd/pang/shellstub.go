// +build !darwin
// +build !freebsd
// +build !linux
// +build !windows

package main

import "log"

func openfolder(folder string) error {
	log.Println(folder)
	return nil
}
