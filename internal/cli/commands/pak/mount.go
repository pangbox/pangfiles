package pak

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/internal/region"
	"github.com/pangbox/pangfiles/internal/shell"
	"github.com/pangbox/pangfiles/pak"
)

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

func (p *cmdPakMount) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	argc := f.NArg()
	argv := f.Args()

	if argc < 2 {
		log.Println("Not enough arguments (did you specify a mount point?)")
		return subcommands.ExitUsageError
	}

	pakfiles := argv[:argc-1]
	mountpoint := argv[argc-1]

	err := os.MkdirAll(mountpoint, 0o775)
	if err != nil {
		log.Printf("Warning: couldn't make mount dir: %v", err)
	}

	fs, err := pak.LoadPaks(region.PakKey(p.region, pakfiles), pakfiles)
	if err != nil {
		log.Fatalf("Loading pak files: %v", err)
	}

	// We don't currently have a good callback for when fuse mounting has succeeded.
	go func() {
		for i := 0; i < 50; i++ {
			time.Sleep(100 * time.Millisecond)
			if stat, err := os.Stat(mountpoint); !os.IsNotExist(err) {
				if stat.IsDir() {
					if err := shell.OpenFolder(mountpoint); err != nil {
						fmt.Printf("Tried to mount folder %s, failed: %v\n", mountpoint, err)
					}
				}
				return
			}
		}
		fmt.Println("Timed out waiting for mount point")
	}()

	if err := fs.Mount(mountpoint); err != nil {
		log.Printf("Mounting filesystem: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func MountCommand() subcommands.Command {
	return &cmdPakMount{}
}
