package config

import "github.com/spf13/viper"

type Redis struct {
	Host     string // Redis.Host
	Username string // Redis.Username
	Password string // Redis.Password
	Database int    // Redis.Database
}

func (r *RealtimeConfig) reloadRedis() {
	r.redis.Host = viper.GetString("Redis.Host")
	r.redis.Username = viper.GetString("Redis.Username")
	r.redis.Password = viper.GetString("Redis.Password")
	r.redis.Database = viper.GetInt("Redis.Database")
}
