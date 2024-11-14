package config

import "github.com/spf13/viper"

type Database struct {
	ConnectionString      string // Database.ConnectionString
	SlaveConnectionString string // Database.SlaveConnectionString
	Dialect               string // Database.Dialect
	LogMode               bool   // Database.LogMode
}

func (r *RealtimeConfig) reloadDatabase() {
	r.database.ConnectionString = viper.GetString("Database.ConnectionString")
	r.database.SlaveConnectionString = viper.GetString("Database.SlaveConnectionString")

	r.database.Dialect = viper.GetString("Database.Dialect")
	r.database.LogMode = viper.GetBool("Database.LogMode")

	r.testDatabase()
}

func (r *RealtimeConfig) testDatabase() {
	testEmptyString(r.database, "ConnectionString")
	testEmptyString(r.database, "Dialect")
}
