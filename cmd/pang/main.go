package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
)

var xteaKeys = []pyxtea.Key{
	pyxtea.KeyUS,
	pyxtea.KeyJP,
	pyxtea.KeyTH,
	pyxtea.KeyEU,
	pyxtea.KeyID,
	pyxtea.KeyKR,
}

var regionToKey = map[string]pyxtea.Key{
	"us": pyxtea.KeyUS,
	"jp": pyxtea.KeyJP,
	"th": pyxtea.KeyTH,
	"eu": pyxtea.KeyEU,
	"id": pyxtea.KeyID,
	"kr": pyxtea.KeyKR,
}

var keyToRegion = map[pyxtea.Key]string{
	pyxtea.KeyUS: "us",
	pyxtea.KeyJP: "jp",
	pyxtea.KeyTH: "th",
	pyxtea.KeyEU: "eu",
	pyxtea.KeyID: "id",
	pyxtea.KeyKR: "kr",
}

func getRegionKey(regionCode string) pyxtea.Key {
	key, ok := regionToKey[regionCode]
	if !ok {
		log.Fatalf("Invalid region %q (valid regions: us, jp, th, eu, id, kr)", regionCode)
	}
	return key
}

func getKeyRegion(key pyxtea.Key) string {
	region, ok := keyToRegion[key]
	if !ok {
		panic("programming error: unexpected key")
	}
	return region
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmdPakMount{}, "")
	subcommands.Register(&cmdPakExtract{}, "")
	subcommands.Register(&cmdUpdateListServe{}, "")
	subcommands.Register(&cmdUpdateListEncrypt{}, "")
	subcommands.Register(&cmdUpdateListDecrypt{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
