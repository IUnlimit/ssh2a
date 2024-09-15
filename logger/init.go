package logger

import (
	"github.com/IUnlimit/ssh2a/conf"
	"github.com/IUnlimit/ssh2a/configs"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

const Prefix = "SSH2A"

func Init() {
	initLog(conf.Config.Log, configs.ParentPath)
}

func initLog(conf *configs.Log, parentPath string) {
	rotateOptions := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour * 24),
	}
	rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(conf.Aging))
	if conf.ForceNew {
		rotateOptions = append(rotateOptions, rotatelogs.ForceNewFile())
	}

	w, err := rotatelogs.New(path.Join(parentPath+"/logs", "%Y-%m-%d.log"), rotateOptions...)
	if err != nil {
		log.Errorf("Rotatelogs init err: %v", err)
		panic(err)
	}

	levels := GetLogLevel(conf.Level)
	log.SetLevel(levels[0]) // hook levels doesn't work
	log.SetReportCaller(true)
	consoleFormatter := LogFormat{
		Prefix:      Prefix,
		EnableColor: conf.Colorful,
	}
	fileFormatter := LogFormat{
		Prefix:      Prefix,
		EnableColor: false,
	}
	Hook = NewLocalHook(w, consoleFormatter, fileFormatter, levels...)
	log.AddHook(Hook)
}
