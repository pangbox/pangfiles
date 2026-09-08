package version

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/version"
)

type versionCmd struct{}

func (versionCmd) Name() string           { return "version" }
func (versionCmd) Synopsis() string       { return "Prints version information to stdout." }
func (v versionCmd) Usage() string        { return fmt.Sprintf("%s:\n  %s\n", v.Name(), v.Synopsis()) }
func (versionCmd) SetFlags(*flag.FlagSet) {}
func (versionCmd) Execute(context.Context, *flag.FlagSet, ...any) subcommands.ExitStatus {
	versionStr := "v" + version.Release
	if version.GitCommit != "" {
		versionStr += "+" + version.GitCommit
	}
	fmt.Println(versionStr)
	return subcommands.ExitSuccess
}

func Command() subcommands.Command {
	return &versionCmd{}
}
