package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/gou/pbe"
	"github.com/bingoohuang/tool/mci"
	"github.com/elliotchance/pie/pie"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const version = "Version: 1.11.9 2020-12-02 10:23:49"

func main() {
	removeSlaves := pflag.StringP("removeSlaves", "", "", "remove slave nodes from cluster, eg 192.168.1.1,192.168.1.2")
	resetLocal := pflag.BoolP("reset", "", false, "reset MySQL cluster")
	readips := pflag.BoolP("readips", "r", false, "read haproxy server ips")

	// --checkmc=checkmc，输出OK表示集群状态是好的，输出其它内容，为详细错误
	checkmc := pflag.StringP("checkmc", "", "", "check mysql cluster, format checkmc/json/table")
	checkmysql := pflag.BoolP("checkmysql", "", false, "check mysql connection")
	ver := pflag.BoolP("version", "v", false, "show version")
	conf := pflag.StringP("config", "c", "./config.toml", "config file path")

	mci.DeclarePflagsByStruct(mci.Settings{})

	pbe.DeclarePflags()

	pflag.Parse()

	checkIllegalArgs()
	printVersion(*ver, version)

	viper.SetEnvPrefix("MCI")
	viper.AutomaticEnv()
	_ = viper.BindPFlags(pflag.CommandLine)

	pbe.DealPflag()

	viper.Set("NoneSetup", *checkmc != "" || *readips || *checkmysql)

	configFile, _ := homedir.Expand(*conf)
	settings := mustLoadConfig(configFile)

	checkSth(settings, *checkmc, *checkmysql, *readips)
	fmt.Println(version)

	removeSlavesFromCluster(*removeSlaves, settings)
	resetLocalMySQLClusterNode(*resetLocal, settings)

	if _, err := settings.CreateMySQLCluster(); err != nil {
		logrus.Errorf("CreateMySQLCluster %v", err)
		os.Exit(1)
	}
}

func printVersion(ver bool, version string) {
	if ver {
		fmt.Println(version)
		os.Exit(0)
	}
}

func checkIllegalArgs() {
	args := pflag.Args()
	if len(args) == 0 {
		return
	}

	fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
	pflag.PrintDefaults()

	os.Exit(1)
}

func checkSth(settings *mci.Settings, checkmc string, checkmysql, readips bool) {
	if checkmc != "" || checkmysql || readips {
		settings.NoLog = true
	}

	if checkmc != "" {
		settings.CheckMySQLCluster(checkmc)
		os.Exit(0)
	}

	if checkmysql {
		settings.CheckMySQL()
		os.Exit(0)
	}

	if readips {
		if mysqlServerAddrs, err := settings.ReadMySQLServersFromHAProxyCfg(); err != nil {
			logrus.Fatal(err)
		} else {
			mysqlServerIPs := pie.Strings(mysqlServerAddrs).Map(func(address string) string {
				pos := strings.LastIndex(address, ":")
				return address[:pos]
			}).Join("\n")
			fmt.Println(mysqlServerIPs)
		}

		os.Exit(0)
	}
}

func removeSlavesFromCluster(removeSlaves string, settings *mci.Settings) {
	if removeSlaves == "" {
		return
	}

	if err := settings.RemoveSlavesFromCluster(removeSlaves); err != nil {
		logrus.Errorf("ResetLocalMySQLClusterNode %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func resetLocalMySQLClusterNode(resetMe bool, settings *mci.Settings) {
	if !resetMe {
		return
	}

	if err := settings.ResetLocalMySQLClusterNode(); err != nil {
		logrus.Errorf("ResetLocalMySQLClusterNode %v", err)
		os.Exit(1)
	}

	os.Exit(0)
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

	// nolint:goerr113
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

func mustLoadConfig(configFile string) *mci.Settings {
	c, _ := loadConfig(configFile)
	mci.ViperToStruct(&c)

	if !viper.GetBool("NoneSetup") {
		c.Setup()
	}

	return &c
}
