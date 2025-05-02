package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"strings"
)

var ()

func main() {
	setupConfig()
	var log = setupLogger()
	log.Infof("Mikrotik Exporter starting...")

	//srv := &http.Server{}
	//if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
	//	logger.Error("Error starting HTTP server", "err", err)
	//	os.Exit(1)
	//}
}

func setupConfig() {
	viper.SetConfigName("mex-settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/mikrotik-exporter")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("MEX")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pflag.String("log.level", "", "Log level (overrides YAML)")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func setupLogger() *logrus.Logger {
	var log = logrus.New()
	logLevel := viper.GetString("log.level")
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Warnf("Invalid log level in config: %s. Using 'info'.", logLevel)
		level = logrus.InfoLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return log
}
