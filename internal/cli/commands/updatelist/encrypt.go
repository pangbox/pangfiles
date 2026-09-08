package updatelist

import (
	"context"
	"flag"
	"log"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/internal/region"
)

type cmdUpdateListEncrypt struct {
	region string
}

func (*cmdUpdateListEncrypt) Name() string { return "updatelist-encrypt" }
func (*cmdUpdateListEncrypt) Synopsis() string {
	return "encrypts an updatelist"
}
func (*cmdUpdateListEncrypt) Usage() string {
	return `updatelist-encrypt [-region <code>] [input file] [output file]:
	Encrypts an updatelist XML document for use with a client.

	When input file is not specified, it defaults to stdin.
	When output file is not specified, it defaults to stdout.

`
}

func (p *cmdUpdateListEncrypt) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.region, "region", "us", "region to use (us, jp, th, eu, id, kr)")
}

func (p *cmdUpdateListEncrypt) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if f.NArg() > 2 {
		log.Println("Too many arguments specified.")
		return subcommands.ExitUsageError
	}
	in := openfile(f.Arg(1))
	out := createfile(f.Arg(2))
	defer closefiles(in, out)
	if err := pyxtea.EncipherStreamPadNull(region.Key(p.region), in, out); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func EncryptCommand() subcommands.Command {
	return &cmdUpdateListEncrypt{}
}
