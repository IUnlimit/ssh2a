package configs

import "embed"

//go:embed config.yml
var ConfigFile embed.FS

const ParentPath = "./"

const ConfigFileName = "config.yml"
