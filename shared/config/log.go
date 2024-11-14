package config

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Log struct {
	Path   string        // Log.Path
	Level  logrus.Level  // Log.Level
	MaxAge time.Duration // Log.MaxAge
}

func (r *RealtimeConfig) reloadLog() {
	r.log.Path = viper.GetString("Log.Path")

	level := viper.GetString("Log.Level")
	if len(level) == 0 {
		fmt.Println("Warning! Empty log level. Set log level to \"DEBUG\".")
		level = "debug"
	}

	var err error
	if r.log.Level, err = logrus.ParseLevel(level); err != nil {
		fmt.Printf("Warning! Invalid log level \"%s\". Set log level to \"DEBUG\" by default.\n", level)
		r.log.Level = logrus.DebugLevel
	}

	maxAge := viper.GetInt("Log.MaxAge")
	if maxAge < 7 {
		fmt.Printf("Warning! Max log age %d less than 7 days. Set max log age to 7 days.\n", maxAge)
		maxAge = 7
	}

	r.log.MaxAge = time.Duration(24*maxAge) * time.Hour

	r.testLog()
}

func (r *RealtimeConfig) testLog() {
	stat, err := os.Stat(r.log.Path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(r.log.Path, os.ModePerm); err != nil {
			panic(fmt.Sprintf("Config - Log.Path. Can not create folder %s", r.log.Path))
		}
	} else if !stat.IsDir() {
		panic(fmt.Sprintf("Config - Log.Path. Path %s not a folder", r.log.Path))
	}

	testEmptyString(r.log, "Path")
}
