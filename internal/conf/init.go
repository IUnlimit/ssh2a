package conf

import (
	"github.com/IUnlimit/ssh2a/configs"
	global "github.com/IUnlimit/ssh2a/internal"
	log "github.com/sirupsen/logrus"
)

func Init() {
	fileFolder := global.ParentPath + "/"
	_, err := LoadConfig(configs.ConfigFileName, fileFolder, "yaml", configs.Config, &global.Config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}
