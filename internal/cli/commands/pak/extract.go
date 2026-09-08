package pak

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/internal/region"
	"github.com/pangbox/pangfiles/pak"
)

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

func (p *cmdPakExtract) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if f.NArg() < 1 {
		log.Println("Not enough arguments. Specify a pak or set of paks to extract.")
		return subcommands.ExitUsageError
	}

	if p.out != "" {
		if err := os.MkdirAll(p.out, 0o775); err != nil {
			log.Printf("Warning: couldn't make output dir: %v", err)
		}
	}

	fs, err := pak.LoadPaks(region.PakKey(p.region, f.Args()), f.Args())
	if err != nil {
		log.Printf("Loading pak files: %v", err)
		return subcommands.ExitFailure
	}

	if p.flat {
		if err = fs.ExtractFlat(p.out); err != nil {
			log.Printf("Extracting pak files: %v", err)
			return subcommands.ExitFailure
		}
	} else {
		if err = fs.Extract(p.out); err != nil {
			log.Printf("Extracting pak files: %v", err)
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess
}

func ExtractCommand() subcommands.Command {
	return &cmdPakExtract{}
}
