package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/tool/mysqlclusterinit"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	conf := pflag.StringP("config", "c", "./config.toml", "config file path")
	pflag.Parse()

	args := pflag.Args()
	if len(args) > 0 {
		fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
		pflag.PrintDefaults()
		os.Exit(0)
	}

	configFile, _ := homedir.Expand(*conf)
	settings := mustLoadConfig(configFile)

	if r := settings.InitMySQLCluster(); r.Error != nil {
		logrus.Panic(r.Error)
	}
}

func loadConfig(configFile string) (config mysqlclusterinit.Settings, err error) {
	if _, err = toml.DecodeFile(configFile, &config); err != nil {
		logrus.Errorf("DecodeFile error %v", err)
	}

	return
}

func mustLoadConfig(configFile string) (config mysqlclusterinit.Settings) {
	var err error
	if config, err = loadConfig(configFile); err != nil {
		logrus.Panic(err)
	}

	if config.Port <= 0 {
		config.Port = 3306
	}

	logrus.Debugf("config: %+v\n", config)
	return config
}
