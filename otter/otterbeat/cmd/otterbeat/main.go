package main

import (
	"fmt"

	"github.com/bingoohuang/gou/cnf"
	"github.com/bingoohuang/gou/enc"
	"github.com/bingoohuang/otterbeat/otter"
	"github.com/spf13/pflag"
)

func main() {
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

	if config.PrintConfig {
		fmt.Printf("Config%s\n", enc.JSONPretty(config))
	}
}
