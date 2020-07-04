package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/pak"
)

var (
	key    pyxtea.Key
	region = flag.String("region", "us", "Region to use (us, jp, th, eu, id, kr)")
)

func init() {
	flag.Parse()
	switch *region {
	case "us":
		key = pyxtea.KeyUS
	case "jp":
		key = pyxtea.KeyJP
	case "th":
		key = pyxtea.KeyTH
	case "eu":
		key = pyxtea.KeyEU
	case "id":
		key = pyxtea.KeyID
	case "kr":
		key = pyxtea.KeyKR
	default:
		log.Fatalf("Invalid region %q (valid regions: us, jp, th, eu, id, kr)", *region)
	}
}

func main() {
	switch flag.Arg(0) {
	case "mount":
		if flag.NArg() < 3 {
			log.Fatalln("Command serve requires 2 argument (pak files/pak glob pattern, mount path)")
		}
		mount(flag.Args()[1:flag.NArg()-1], flag.Arg(flag.NArg()-1))
	default:
		log.Fatalln("Please provide a valid command. (valid commands: mount)")
	}
}

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
