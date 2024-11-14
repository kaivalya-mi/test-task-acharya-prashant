package main

import (
	"context"
	"flag"
	"os"
	"test-task/model"
	"test-task/routers/api"
	"test-task/shared/cache"
	"test-task/shared/config"
	"test-task/shared/database"
	"test-task/shared/log"
	"test-task/shared/utils"
)

// Execution starts from main function
func main() {
	configFile := flag.String("c", "", "configuration file without extension. For config.toml then put \" -c config-development\"")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	var cf config.IConfig
	if *configFile == "" {
		cf = config.NewRealtimeConfig("config", pwd)
	} else {
		cf = config.NewConfig(*configFile)
	}

	log.Init("app", cf.AppVersion(), cf.Info().Path, cf.Info().Level, cf.Info().MaxAge)
	log.GetLog().Info("", "App service start! %s", cf.AppVersion())
	database.Init(cf)
	cache.CreateConnection(cf)

	log.GetLog().Info("", "DB connected")
	model.AutoMigrate()
	rt := api.NewRouter(cf)
	rt.Setup()

	go rt.Run()

	utils.GracefulStop(log.GetLog(), func(ctx context.Context) error {
		var err error
		if err = rt.Close(ctx); err != nil {
			return err
		}
		if err = database.Close(); err != nil {
			return err
		}
		return nil
	})
}
