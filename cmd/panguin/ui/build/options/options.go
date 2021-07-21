package options

import "github.com/evanw/esbuild/pkg/api"

var BuildOptions = api.BuildOptions{
	Bundle:            true,
	MinifySyntax:      true,
	MinifyWhitespace:  true,
	MinifyIdentifiers: true,
	Outfile:           "./ui/dist/index.js",
	EntryPoints:       []string{"./ui/src/index.tsx"},
	Write:             true,
	AllowOverwrite:    true,
}
