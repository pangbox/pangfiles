// +build devmode

package main

import (
	"log"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/jchv/go-webview-selector"
	"github.com/pangbox/pangfiles/cmd/panguin/ui/build/options"
)

func run(setup func(webview.WebView)) {
	w := webview.New(true)
	defer w.Destroy()

	setup(w)

	reload := func(b api.BuildResult) {
		log.Println("build finished")
		for _, m := range b.Errors {
			log.Println("error:", m)
		}
		for _, m := range b.Warnings {
			log.Println("warning:", m)
		}

		w.Dispatch(func() {
			w.Navigate(makepage(b.OutputFiles[0].Contents, b.OutputFiles[1].Contents))
		})
	}

	options.BuildOptions.Watch = &api.WatchMode{
		OnRebuild: reload,
	}

	result := api.Build(options.BuildOptions)
	reload(result)

	w.Run()
}
