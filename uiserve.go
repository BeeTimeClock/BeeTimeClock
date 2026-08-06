//go:build !noui
// +build !noui

package main

import "embed"

//go:embed ui/dist/spa
var uiFS embed.FS

func init() {
	IsUiEmbedded = true
}
