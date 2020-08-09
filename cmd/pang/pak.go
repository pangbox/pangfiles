package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/pak"
)

func getPakKey(region string, patterns []string) pyxtea.Key {
	var key pyxtea.Key
	if region == "" {
		log.Println("Auto-detecting pak region (use -region to improve startup delay.)")
		key = pak.MustDetectRegion(patterns, xteaKeys)
		log.Printf("Detected pak region as %s.", strings.ToUpper(getKeyRegion(key)))
	} else {
		key = getRegionKey(region)
	}
	return key
}

type cmdPakMount struct {
	region string
	flat   bool
	open   bool
}

func (*cmdPakMount) Name() string     { return "pak-mount" }
func (*cmdPakMount) Synopsis() string { return "mounts a set of pak files" }
func (*cmdPakMount) Usage() string {
	return `pak-mount [-flat] [-region <code>] <pak files> <mount point>:
	Mounts a set of ordered pak files as a unified filesystem.
	You can specify globs like projectg*.pak to get PangYa-like behavior.

	On Windows, the mount point must be a drive letter specification, e.g. P:
	On other OSes, the mount point should be a directory, like $HOME/pak.

`
}

func (p *cmdPakMount) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.flat, "flat", false, "flatten the hierarchy (not implemented yet)")
	f.StringVar(&p.region, "region", "", "region to use (us, jp, th, eu, id, kr)")
	f.BoolVar(&p.open, "open", true, "when true (default) open folder upon mounting")
}

func (p *cmdPakMount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	argc := f.NArg()
	argv := f.Args()

	if argc < 2 {
		log.Println("Not enough arguments (did you specify a mount point?)")
		return subcommands.ExitUsageError
	}

	pakfiles := argv[:argc-1]
	mountpoint := argv[argc-1]

	os.MkdirAll(mountpoint, 0o775)

	fs, err := pak.LoadPaks(getPakKey(p.region, f.Args()), pakfiles)
	if err != nil {
		log.Fatal(err)
	}

	// We don't currently have a good callback for when fuse mounting has succeeded.
	go func() {
		for i := 0; i < 50; i++ {
			time.Sleep(100 * time.Millisecond)
			if stat, err := os.Stat(mountpoint); !os.IsNotExist(err) {
				if stat.IsDir() {
					openfolder(mountpoint)
				}
				return
			}
		}
		fmt.Println("Timed out waiting for mount point")
	}()

	if err := fs.Mount(mountpoint); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

type cmdPakExtract struct {
	out    string
	region string
	flat   bool
}

func (*cmdPakExtract) Name() string     { return "pak-extract" }
func (*cmdPakExtract) Synopsis() string { return "extracts a set of pak files" }
func (*cmdPakExtract) Usage() string {
	return `pak-extract [-flat] [-region <code>] [-o <output directory>] <pak files>:
	Extracts a set of pak files into a directory.
	
	This will treat the set of pak files as a single incremental archive.

`
}

func (p *cmdPakExtract) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.out, "o", "", "destination to extract to")
	f.BoolVar(&p.flat, "flat", false, "flatten the hierarchy (not implemented yet)")
	f.StringVar(&p.region, "region", "", "region to use (us, jp, th, eu, id, kr)")
}

func (p *cmdPakExtract) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 {
		log.Println("Not enough arguments. Specify a pak or set of paks to extract.")
		return subcommands.ExitUsageError
	}

	if p.out != "" {
		os.MkdirAll(p.out, 0o775)
	}

	fs, err := pak.LoadPaks(getPakKey(p.region, f.Args()), f.Args())
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	if err = fs.Extract(p.out); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
