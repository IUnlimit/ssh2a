package configs

import "time"

type Config struct {
	Log           *Log   `yaml:"log"`
	RootDir       string `yaml:"root-dir"`
	Authorization string `yaml:"authorization"`
}

type Log struct {
	ForceNew bool          `yaml:"force-new,omitempty"`
	Level    string        `yaml:"level,omitempty"`
	Aging    time.Duration `yaml:"aging,omitempty"`
	Colorful bool          `yaml:"colorful,omitempty"`
}
