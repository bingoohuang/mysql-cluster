package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/gossh/pbe"
	"github.com/bingoohuang/tool/mci"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	readips := pflag.BoolP("readips", "r", false, "read haproxy server ips")
	checkmc := pflag.BoolP("checkmc", "m", false, "check mysql cluster")
	checkmysql := pflag.BoolP("checkmysql", "", false, "check mysql connection")
	ver := pflag.BoolP("version", "v", false, "show version")
	conf := pflag.StringP("config", "c", "./config.toml", "config file path")

	mci.DeclarePflagsByStruct(mci.Settings{})

	pbe.DeclarePflags()

	pflag.Parse()

	args := pflag.Args()
	if len(args) > 0 {
		fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
		pflag.PrintDefaults()
		os.Exit(1)
	}

	if *ver {
		fmt.Printf("Version: 1.5.3\n")
		return
	}

	viper.SetEnvPrefix("MCI")
	viper.AutomaticEnv()
	_ = viper.BindPFlags(pflag.CommandLine)

	pbe.DealPflag()

	configFile, _ := homedir.Expand(*conf)
	settings := mustLoadConfig(configFile)

	if *checkmc {
		settings.CheckMySQLCluster()
	}

	if *checkmysql {
		settings.CheckMySQL()
	}

	if *readips {
		settings.CheckHAProxyServers()
	}

	if *checkmc || *readips || *checkmysql {
		return
	}

	if _, err := settings.InitMySQLCluster(); err != nil {
		logrus.Errorf("error %v", err)
		os.Exit(1)
	}
}

func findConfigFile(configFile string) (string, error) {
	if mci.FileExists(configFile) == nil {
		return configFile, nil
	}

	if ex, err := os.Executable(); err == nil {
		exPath := filepath.Dir(ex)
		configFile = filepath.Join(exPath, "config.toml")
	}

	if mci.FileExists(configFile) == nil {
		return configFile, nil
	}

	return "", errors.New("unable to find config file")
}

func loadConfig(configFile string) (config mci.Settings, err error) {
	if file, err := findConfigFile(configFile); err != nil {
		return config, err
	} else if _, err = toml.DecodeFile(file, &config); err != nil {
		logrus.Errorf("DecodeFile error %v", err)
	}

	return config, err
}

func mustLoadConfig(configFile string) (config mci.Settings) {
	config, _ = loadConfig(configFile)
	mci.ViperToStruct(&config)

	return config
}
