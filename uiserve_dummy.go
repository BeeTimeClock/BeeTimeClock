//go:build noui
// +build noui

package main

import "embed"

var uiFS embed.FS

func init() {
	IsUiEmbedded = false
}
