package updatelistsrv

import (
	"context"
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/internal/region"
)

type cmdUpdateListServe struct {
	region string
	listen string
}

func (*cmdUpdateListServe) Name() string { return "updatelist-serve" }
func (*cmdUpdateListServe) Synopsis() string {
	return "serves an updatelist for a game folder"
}
func (*cmdUpdateListServe) Usage() string {
	return `pak-extract [-region <code>] [-listen <address>] <game folder>:
	Serves an automatically updating updatelist for a game folder.

`
}

func (p *cmdUpdateListServe) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.region, "region", "", "region to use (us, jp, th, eu, id, kr)")
	f.StringVar(&p.listen, "listen", ":8080", "address to listen on")
}

func (p *cmdUpdateListServe) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if f.NArg() > 1 {
		log.Println("Too many arguments.")
		return subcommands.ExitUsageError
	} else if f.NArg() < 1 {
		log.Println("Not enough arguments. Try specifying a game folder.")
		return subcommands.ExitUsageError
	}

	dir := f.Arg(0)

	key := region.PakKey(p.region, []string{
		filepath.Join(dir, "projectg*.pak"),
		filepath.Join(dir, "ProjectG*.pak"),
	})

	s := server{
		key:   key,
		dir:   dir,
		cache: map[string]cacheentry{},
	}
	if err := http.ListenAndServe(p.listen, &s); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func ServeCommand() subcommands.Command {
	return &cmdUpdateListServe{}
}
