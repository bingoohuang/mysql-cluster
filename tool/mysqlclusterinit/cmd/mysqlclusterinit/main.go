package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/tool/mysqlclusterinit"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	checkmysql := pflag.BoolP("checkmysql", "m", false, "check mysql")
	ver := pflag.BoolP("version", "v", false, "show version")
	conf := pflag.StringP("config", "c", "./config.toml", "config file path")

	mysqlclusterinit.DeclarePflagsByStruct(mysqlclusterinit.Settings{})

	pflag.Parse()

	args := pflag.Args()
	if len(args) > 0 {
		fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
		pflag.PrintDefaults()
		os.Exit(1)
	}

	if *ver {
		fmt.Printf("Version: 1.3.1\n")
		return
	}

	viper.SetEnvPrefix("MCI")
	viper.AutomaticEnv()
	_ = viper.BindPFlags(pflag.CommandLine)

	configFile, _ := homedir.Expand(*conf)
	settings := mustLoadConfig(configFile)

	if *checkmysql {
		settings.CheckMySQL()
		return
	}

	if _, err := settings.InitMySQLCluster(); err != nil {
		logrus.Errorf("error %v", err)
		os.Exit(1)
	}
}

func loadConfig(configFile string) (config mysqlclusterinit.Settings, err error) {
	if mysqlclusterinit.FileExists(configFile) != nil {
		return config, nil
	}

	if _, err = toml.DecodeFile(configFile, &config); err != nil {
		logrus.Errorf("DecodeFile error %v", err)
	}

	return
}

func mustLoadConfig(configFile string) (config mysqlclusterinit.Settings) {
	config, _ = loadConfig(configFile)
	mysqlclusterinit.ViperToStruct(&config)

	return config
}
