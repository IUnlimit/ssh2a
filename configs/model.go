package configs

import "time"

type Config struct {
	Log           *Log           `yaml:"log"`
	Bind          *Bind          `yaml:"bind,omitempty"`
	Authorization *Authorization `yaml:"authorization"`
	Database      *Database      `yaml:"database"`
	Honeypot      *Honeypot      `yaml:"honeypot"`
	Auth          *AuthPolicy    `yaml:"auth"`
	Admin         *Admin         `yaml:"admin"`
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

type Database struct {
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	DBName   string `yaml:"dbname,omitempty"`
}

type Honeypot struct {
	TriggerTimeout time.Duration `yaml:"trigger-timeout,omitempty"`
}

type AuthPolicy struct {
	ValidDuration time.Duration `yaml:"valid-duration,omitempty"`
}

type Admin struct {
	AllowedHosts []string `yaml:"allowed-hosts,omitempty"`
}
