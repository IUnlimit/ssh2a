package configs

import "time"

type Config struct {
	Log           *Log           `yaml:"log"`
	Bind          *Bind          `yaml:"bind,omitempty"`
	Authorization *Authorization `yaml:"authorization"`
}

type Log struct {
	ForceNew bool          `yaml:"force-new,omitempty"`
	Level    string        `yaml:"level,omitempty"`
	Aging    time.Duration `yaml:"aging,omitempty"`
	Colorful bool          `yaml:"colorful,omitempty"`
}

type Bind struct {
	Host     string `yaml:"host,omitempty"`
	HttpPort int    `yaml:"http-port,omitempty"`
	SSHPort  int    `yaml:"ssh-port,omitempty"`
}

type Authorization struct {
	Type          string         `yaml:"type,omitempty"`
	Basic         *Basic         `yaml:"basic"`
	Authenticator *Authenticator `yaml:"authenticator"`
}

type Basic struct {
	Secret string `yaml:"secret,omitempty"`
}

type Authenticator struct {
	PrivateSecret string `yaml:"private-secret,omitempty"`
}
