// +build !devmode

package main

import (
	"github.com/jchv/go-webview-selector"
	"github.com/pangbox/pangfiles/cmd/panguin/ui"
)

func run(setup func(webview.WebView)) {
	w := webview.New(false)
	defer w.Destroy()
	setup(w)
	w.Navigate(makepage(ui.Script, ui.Style))
	w.Run()
}
