package ui

import _ "embed"

//go:embed dist/index.js
var Script []byte

//go:embed dist/index.css
var Style []byte
