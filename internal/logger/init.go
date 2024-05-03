package logger

import (
	"path"
	"time"

	global "github.com/IUnlimit/ssh2a/internal"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

const Prefix = "SSH2A"

func Init() {
	initLog()
}

func initLog() {
	config := global.Config
	rotateOptions := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour * 24),
	}
	rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(config.Log.Aging))
	if config.Log.ForceNew {
		rotateOptions = append(rotateOptions, rotatelogs.ForceNewFile())
	}

	w, err := rotatelogs.New(path.Join(global.ParentPath+"/logs", "%Y-%m-%d.log"), rotateOptions...)
	if err != nil {
		log.Errorf("Rotatelogs init err: %v", err)
		panic(err)
	}

	levels := GetLogLevel(config.Log.Level)
	log.SetLevel(levels[0]) // hook levels doesn't work
	log.SetReportCaller(true)
	consoleFormatter := LogFormat{
		Prefix:      Prefix,
		EnableColor: config.Log.Colorful,
	}
	fileFormatter := LogFormat{
		Prefix:      Prefix,
		EnableColor: false,
	}
	Hook = NewLocalHook(w, consoleFormatter, fileFormatter, levels...)
	log.AddHook(Hook)
}
