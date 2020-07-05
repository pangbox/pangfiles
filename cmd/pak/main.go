package main

import (
	"flag"
	"log"
	"os"

	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/pak"
)

var (
	key    pyxtea.Key
	region = flag.String("region", "us", "Region to use (us, jp, th, eu, id, kr)")
	flat   = flag.Bool("flat", false, "Flatten the directory structure to match PangYa's internal view.")
	dir    = flag.String("dir", "", "Directory to mount to or extract to.")
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
		if flag.NArg() < 2 {
			log.Fatalln("Command serve requires a glob or list of pak files.")
		}
		if *dir == "" {
			log.Fatalln("-dir not specified (specify -dir $PWD if you really want to mount into current directory)")
		}
		os.MkdirAll(*dir, 0o775)
		mount(flag.Args()[1:flag.NArg()], *dir)
	case "extract":
		if flag.NArg() < 2 {
			log.Fatalln("Command extract requires a glob or list of pak files.")
		}
		if *dir == "" {
			log.Fatalln("-dir not specified (specify -dir $PWD if you really want to extract into current directory)")
		}
		extract(flag.Args()[1:flag.NArg()], *dir)
	default:
		log.Fatalln("Please provide a valid command. (valid commands: mount, extract)")
	}
}

func mount(patterns []string, mountpoint string) {
	fs, err := pak.LoadPaks(key, patterns)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(fs.Mount(mountpoint))
}

func extract(patterns []string, dest string) {
	fs, err := pak.LoadPaks(key, patterns)
	if err != nil {
		log.Fatal(err)
	}

	err = fs.Extract(dest)
	if err != nil {
		log.Fatal(err)
	}
}
