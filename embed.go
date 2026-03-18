package main

import "embed"

//go:embed all:web/dist
var webDist embed.FS
