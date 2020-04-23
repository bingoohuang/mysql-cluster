package main

import (
	"github.com/bingoohuang/gou/sy"
	"github.com/bingoohuang/otterbeat/otter"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config := otter.Config{}

	sy.SetupApp(&sy.AppOption{
		EnvPrefix:   "OTTERBEAT",
		LogLevel:    "debug",
		ConfigBeans: []interface{}{&config},
	})

	config.Run()
}
