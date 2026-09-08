package updatelist

import (
	"context"
	"flag"
	"log"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/internal/region"
)

type cmdUpdateListDecrypt struct {
	region string
}

func (*cmdUpdateListDecrypt) Name() string { return "updatelist-decrypt" }
func (*cmdUpdateListDecrypt) Synopsis() string {
	return "decrypts an updatelist"
}
func (*cmdUpdateListDecrypt) Usage() string {
	return `updatelist-decrypt [-region <code>] [input file] [output file]:
	Decrypts an encrypted updatelist back to plaintext XML.

	When input file is not specified, it defaults to stdin.
	When output file is not specified, it defaults to stdout.

`
}

func (p *cmdUpdateListDecrypt) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.region, "region", "us", "region to use (us, jp, th, eu, id, kr)")
}

func (p *cmdUpdateListDecrypt) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if f.NArg() > 2 {
		log.Println("Too many arguments specified.")
		return subcommands.ExitUsageError
	}
	in := openfile(f.Arg(0))
	out := createfile(f.Arg(1))
	defer closefiles(in, out)
	if err := pyxtea.DecipherStreamTrimNull(region.Key(p.region), in, out); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func DecryptCommand() subcommands.Command {
	return &cmdUpdateListDecrypt{}
}
