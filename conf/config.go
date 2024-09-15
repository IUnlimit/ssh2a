package conf

import (
	"embed"
	"github.com/IUnlimit/ssh2a/configs"
	"github.com/IUnlimit/ssh2a/tools"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"os"
)

var Config configs.Config

func init() {
	err := tryCreateConfigFile(configs.ParentPath, configs.ConfigFileName, configs.ConfigFile)
	if err != nil {
		log.Errorf("Try to create config file failed, %v", err)
	}
	err = configor.New(&configor.Config{
		AutoReload: true,
		AutoReloadCallback: func(config interface{}) {
			Config = config.(configs.Config)
			log.Infof("ConfigFile\n`%s`\nchanged, auto reloading ...", spew.Sdump(config))
		},
	}).Load(&Config, configs.ConfigFileName)
	if err != nil {
		log.Fatalf("Configuration load failed, %v", err)
	}
}

func tryCreateConfigFile(fileFolder string, fileName string, fs embed.FS) error {
	filePath := fileFolder + fileName
	exists := tools.FileExists(filePath)
	if !exists {
		log.Warnf("Can't find `%s`, generating default configuration", fileName)
		data, err := fs.ReadFile(fileName)
		if err != nil {
			return err
		}
		err = os.MkdirAll(fileFolder, os.ModePerm)
		if err != nil {
			return err
		}
		err = tools.CreateFile(filePath, data)
		if err != nil {
			return err
		}
	}
	return nil
}
