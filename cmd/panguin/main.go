package main

import (
	"github.com/jchv/go-webview-selector"
	"github.com/vincent-petithory/dataurl"
)

//go:generate go run ./ui/build

func main() {
	run(func(w webview.WebView) {
		w.SetTitle("Panguin")
		w.SetSize(800, 600, webview.HintNone)
	})
}

func makepage(script, style []byte) string {
	return "data:text/html,<script>" + dataurl.Escape(script) + "</script>" + "<style>" + dataurl.Escape(style) + "</style>"
}
