package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"mikrotik-exporter/internal/collector"
	"mikrotik-exporter/internal/logger"

	"strings"
)

func main() {
	setupConfig()
	logger.Setup()
	logger.Log.Infof("Mikrotik Exporter starting...")

	collector.Start()
}

func setupConfig() {
	viper.SetConfigName("mex-settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/mikrotik-exporter")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("MEX")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pflag.String("log.level", "", "Log level (overrides mex-settings.yaml file)")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
