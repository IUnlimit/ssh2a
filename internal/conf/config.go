package conf

import (
	"embed"
	"encoding/json"
	"errors"
	"github.com/IUnlimit/ssh2a/tools"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// LoadConfig creat and load config, return exists(file)
// kind: json / yaml
func LoadConfig(fileName string, fileFolder string, kind string, fs embed.FS, config any) (bool, error) {
	filePath := fileFolder + fileName
	exists := tools.FileExists(filePath)
	if !exists {
		log.Warnf("Can't find `%s`, generating default configuration", fileName)
		data, err := fs.ReadFile(fileName)
		if err != nil {
			return false, err
		}
		err = os.MkdirAll(fileFolder, os.ModePerm)
		if err != nil {
			return false, err
		}
		err = tools.CreateFile(filePath, data)
		if err != nil {
			return false, err
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return exists, err
	}

	if kind == "json" {
		err = json.Unmarshal(data, config)
	} else if kind == "yaml" {
		err = yaml.Unmarshal(data, config)
	} else {
		err = errors.New("unknown file type: " + kind)
	}
	if err != nil {
		return exists, err
	}
	return exists, nil
}
