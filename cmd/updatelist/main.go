package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/util"
)

var (
	key    pyxtea.Key
	region = flag.String("region", "us", "Region to use (us, jp, th, eu, id, kr)")
	listen = flag.String("listen", ":8080", "Address to listen on.")
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
	case "serve":
		if flag.NArg() != 2 {
			log.Fatalln("Command serve requires 1 argument (path to game folder)")
		}
		serve(*listen, flag.Arg(1), key)
	case "encrypt", "decrypt":
		if flag.NArg() == 1 {
			log.Println("Reading from stdin.")
		} else if flag.NArg() > 3 {
			log.Fatalf("Command %s requires at most 2 arguments (input file, output file)", flag.Arg(0))
		}
		encrypt := flag.Arg(0) == "encrypt"
		in := openfile(flag.Arg(1))
		out := createfile(flag.Arg(2))
		defer closefiles(in, out)
		crypt(in, out, key, encrypt)
	default:
		log.Fatalln("Please provide a valid command. (valid commands: serve, encrypt, decrypt)")
	}
}

func serve(listen, dir string, key pyxtea.Key) {
	s := server{
		key:   key,
		dir:   dir,
		cache: map[string]cacheentry{},
	}
	log.Fatalln(http.ListenAndServe(listen, &s))
}

func crypt(in io.Reader, out io.Writer, key pyxtea.Key, encrypt bool) {
	if encrypt {
		pyxtea.EncipherStream(key, util.NullInputPadder{Reader: in}, out)
	} else {
		pyxtea.DecipherStream(key, in, out)
	}
}

func openfile(infile string) *os.File {
	if infile == "" {
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
