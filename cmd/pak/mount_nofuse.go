// +build !freebsd,!linux

package main

import "log"

func mount(patterns []string, mountpoint string) {
	log.Fatalln("Sorry, pak fuse is not supported on this platform yet. :(")
}
