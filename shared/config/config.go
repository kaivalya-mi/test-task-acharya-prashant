package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"test-task/shared/utils"

	"github.com/spf13/viper"
)

// IConfig is
type IConfig interface {
	AppVersion() string
	AppRegion() string

	Env() string
	WorkDir() string

	App() *App
	Info() *Log

	Database() *Database
	Redis() *Redis
}

// RealtimeConfig is
type RealtimeConfig struct {
	appVersion string
	appRegion  string

	env     string // Server.ENV
	workDir string // Server.WorkDir

	log      Log
	redis    Redis
	app      App
	database Database
}

func testEmptyString(entity interface{}, path string) {
	s := utils.GetStringFromStruct(entity, path)
	if strings.TrimSpace(s) == "" {
		name := utils.GetStructName(entity)
		panic(fmt.Sprintf("Config - %s.%s can not be empty", name, path))
	}
}

// NewRealtimeConfig is
func NewRealtimeConfig(configName, path string) IConfig {
	viper.SetConfigName(configName)

	// TODO move it to the same folder as executable apps stay
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	rc := RealtimeConfig{}

	rc.reload()

	return &rc
}

// NewRealtimeConfig is
func NewConfig(fullPath string) IConfig {
	basedir := filepath.Dir(fullPath)

	filename := filepath.Base(fullPath)
	extname := filepath.Ext(filename)
	if len(extname) > 0 {
		filename = filename[:len(filename)-len(extname)]
	}

	return NewRealtimeConfig(filename, basedir)
}

func (r *RealtimeConfig) reload() {
	r.appVersion = os.Getenv("APP_VERSION")
	r.appRegion = os.Getenv("APP_REGION")

	r.env = viper.GetString("Server.ENV")

	r.workDir = viper.GetString("Server.WorkDir")
	if len(r.workDir) == 0 {
		r.workDir = "."
	}

	r.reloadLog()
	r.reloadApp()
	r.reloadDatabase()
	r.reloadRedis()
}

func (r *RealtimeConfig) AppVersion() string {
	return r.appVersion
}

func (r *RealtimeConfig) AppRegion() string {
	return r.appRegion
}

func (r *RealtimeConfig) Env() string {
	return r.env
}

func (r *RealtimeConfig) WorkDir() string {
	return r.workDir
}

func (r *RealtimeConfig) Info() *Log {
	return &r.log
}

func (r *RealtimeConfig) App() *App {
	return &r.app
}

func (r *RealtimeConfig) Database() *Database {
	return &r.database
}

func (r *RealtimeConfig) Redis() *Redis {
	return &r.redis
}
