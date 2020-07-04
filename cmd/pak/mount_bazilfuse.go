// +build freebsd linux

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/pangbox/pangfiles/pak"
)

func mount(patterns []string, mountpoint string) {
	pakfs := pak.NewFS(key)
	for _, pattern := range patterns {
		err := pakfs.LoadPaksFromGlob(pattern)
		if err != nil {
			log.Fatalf("Error loading pak files: %s", err)
		}
	}

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("pakfs"),
		fuse.Subtype("pyfs"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	i := make(chan os.Signal, 1)
	signal.Notify(i, os.Interrupt)
	go func() {
		<-i
		fmt.Println("Received interrupt, exiting.")
		fuse.Unmount(mountpoint)
		os.Exit(0)
	}()

	err = fs.Serve(c, pakfs)
	if err != nil {
		log.Fatal(err)
	}
}
