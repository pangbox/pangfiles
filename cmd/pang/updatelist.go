package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
)

func openfile(infile string) *os.File {
	if infile == "" {
		log.Println("Reading from stdin.")
		return os.Stdin
	}
	in, err := os.Open(infile)
	if err != nil {
		log.Fatalf("Error opening input file %q: %s", infile, err)
	}
	return in
}

func createfile(outfile string) *os.File {
	if outfile == "" {
		return os.Stdout
	}
	out, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("Error opening output file %q: %s", outfile, err)
	}
	return out
}

func closefiles(in *os.File, out *os.File) {
	if err := in.Close(); err != nil {
		log.Printf("Warning: an error occurred during close of input file: %s", err)
	}
	if err := out.Close(); err != nil {
		log.Printf("Warning: an error occurred during close of output file: %s", err)
	}
}

type cmdUpdateListServe struct {
	region string
	listen string
}

func (*cmdUpdateListServe) Name() string { return "updatelist-serve" }
func (*cmdUpdateListServe) Synopsis() string {
	return "serves an updatelist for a game folder"
}
func (*cmdUpdateListServe) Usage() string {
	return `pak-extract [-region <code>] [-listen <address>] <pak files>:
	Serves an automatically updating updatelist for a game folder.

`
}

func (p *cmdUpdateListServe) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.region, "region", "us", "region to use (us, jp, th, eu, id, kr)")
	f.StringVar(&p.listen, "listen", ":8080", "address to listen on")
}

func (p *cmdUpdateListServe) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() > 1 {
		log.Println("Too many arguments.")
		return subcommands.ExitUsageError
	} else if f.NArg() < 1 {
		log.Println("Not enough arguments. Try specifying a game folder.")
		return subcommands.ExitUsageError
	}
	s := server{
		key:   getRegionKey(p.region),
		dir:   f.Arg(0),
		cache: map[string]cacheentry{},
	}
	if err := http.ListenAndServe(p.listen, &s); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

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

func (p *cmdUpdateListEncrypt) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() > 2 {
		log.Println("Too many arguments specified.")
		return subcommands.ExitUsageError
	}
	in := openfile(f.Arg(1))
	out := createfile(f.Arg(2))
	defer closefiles(in, out)
	if err := pyxtea.EncipherStreamPadNull(getRegionKey(p.region), in, out); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

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

func (p *cmdUpdateListDecrypt) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() > 2 {
		log.Println("Too many arguments specified.")
		return subcommands.ExitUsageError
	}
	in := openfile(f.Arg(1))
	out := createfile(f.Arg(2))
	defer closefiles(in, out)
	if err := pyxtea.DecipherStreamTrimNull(getRegionKey(p.region), in, out); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
