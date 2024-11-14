package config

import (
	"github.com/spf13/viper"
)

type App struct {
	Port            string
	AccessTokenKey  string
	RefreshTokenKey string
}

func (r *RealtimeConfig) reloadApp() {
	r.app.Port = viper.GetString("App.Port")
	r.app.AccessTokenKey = viper.GetString("App.AccessTokenKey")
	r.app.RefreshTokenKey = viper.GetString("App.RefreshTokenKey")
	r.testApp()
}

func (r *RealtimeConfig) testApp() {
	testEmptyString(r.app, "Port")
}
