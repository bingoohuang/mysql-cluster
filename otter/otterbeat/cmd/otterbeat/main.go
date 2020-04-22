package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bingoohuang/gou/cnf"
	"github.com/bingoohuang/gou/enc"
	"github.com/bingoohuang/otterbeat/otter"
	"github.com/spf13/pflag"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	ver := pflag.BoolP("version", "v", false, "show version")

	cnf.DeclarePflags()
	cnf.DeclarePflagsByStruct(otter.Config{})

	if err := cnf.ParsePflags("OTTERBEAT"); err != nil {
		panic(err)
	}

	if *ver {
		fmt.Println("Version: v0.0.1")
		return
	}

	var config otter.Config

	cnf.LoadByPflag(&config)

	config.SetDefault()

	if config.PrintConfig {
		fmt.Printf("Config%s\n", enc.JSONPretty(config))
	}

	config.Run()
}
