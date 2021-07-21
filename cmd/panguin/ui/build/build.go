package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/pangbox/pangfiles/cmd/panguin/ui/build/options"
)

func main() {
	results := api.Build(options.BuildOptions)
	for _, m := range results.Errors {
		log.Println("error:", m)
	}
	for _, m := range results.Warnings {
		log.Println("warning:", m)
	}
	if len(results.Errors) > 0 {
		os.Exit(1)
	}
}
